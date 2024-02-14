package event

import (
	"bytes"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}
	return consumer, err
}

func (c *Consumer) setup() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(ch)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := declareQueue(ch)
	if err != nil {
		return err
	}
	for _, val := range topics {
		err = ch.QueueBind(
			q.Name,
			val,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	forever := make(chan bool)

	go func() {
		for dat := range messages {
			var payload Payload
			_ = json.Unmarshal(dat.Body, &payload)

			go handlePayload(payload)

		}
	}()
	log.Printf("Waiting for messages [Exchange, Queue] [logs_topic, %s]", q.Name)
	<-forever
	return nil
}

func handlePayload(dat Payload) {
	switch dat.Name {
	case "log":
		err := logEvent(dat)
		if err != nil {
			log.Println("handle payload error log ", err)
		}
	//TODO реализовать все остальные кейс auth mail и тд
	default:
		err := logEvent(dat)
		if err != nil {
			log.Println("handle payload error default ", err)
		}
	}

}

func logEvent(dat Payload) error {
	jsonData, err := json.Marshal(dat)
	if err != nil {
		log.Println("error marshaling log payload from request")
		return err
	}
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
