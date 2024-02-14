package event

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch)
}

func (e *Emitter) Push(ctx context.Context, event, severity string) error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Println("Pushing events onto the queue")

	err = ch.PublishWithContext(ctx,
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		}) //todo сделать потом вместо log_topic возможность создават и другие
	if err != nil {
		return err
	}
	return nil
}
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{connection: conn}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
