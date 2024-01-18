package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/bruce-mig/go-micro/listener/lib/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create a new consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// consumer.Listen watches the queue and consumes events for all the provided topics.
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

// connect tries to connect to RabbitMQ, and delays between attempts.
// If we can't connect after 10 tries (with increasing delays), return an error
func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	var rabbitURL = os.Getenv("RABBIT_URL")

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			// we have a connection to rabbitmq, so set connection = c and break out of
			// this loop
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 10 {
			fmt.Println(err)
			return nil, err
		}

		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), 2)))
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
