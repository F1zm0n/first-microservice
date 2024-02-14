package main

import (
	"github.com/f1zm0n/listener/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"os"
	"time"
)

const rabbitMQURL = "amqp://guest:guest@rabbitmq"

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("listening and consuming for rabbitmq messages")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Panicf("error creating new rabbitmq consumer %v", err)
	}
	err = consumer.Listen([]string{"log.INFO", "log.WARN", "log.ERROR"}) //todo доделать что бы все принимало
	if err != nil {
		log.Println(err)
	}

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var connection *amqp.Connection
	backOff := 1 * time.Second

	for {
		c, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			log.Println("RabbitMQ not ready yet")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			log.Println("couldn't connect to rabbitMQ time limit reached ", err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("rabbitmq connection is backing off")
		time.Sleep(backOff)
	}
	log.Println("connected to rabbitmq")
	return connection, nil
}
