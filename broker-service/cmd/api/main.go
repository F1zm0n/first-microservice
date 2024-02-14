package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const port = "8080"
const rabbitMQURL = "amqp://guest:guest@rabbitmq"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{Rabbit: rabbitConn}

	log.Printf("starting broker service on port %s\n", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panicf("error starting broker service on port %s, err: %v", port, err)
	}
	log.Printf("broker server is running on port %s", port)
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
