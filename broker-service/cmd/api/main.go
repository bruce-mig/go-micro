package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// webPort the port that we listen on for api calls
const webPort = "80"

// Config is the type we'll use as a receiver to share application
// configuration around our app.
type Config struct {
	Rabbit         *amqp.Connection
	LogServiceURLs map[string]string
	//MailServiceURLs map[string]string
	//AuthServiceURLs map[string]string
}

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connectToRabbit()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// connectToRabbit tries to connect to RabbitMQ, for up to 30 seconds
func connectToRabbit() (*amqp.Connection, error) {
	var rabbitConn *amqp.Connection
	var counts int64
	var rabbitURL = os.Getenv("RABBIT_URL")

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			rabbitConn = c
			break
		}

		if counts > 15 {
			fmt.Println(err)
			return nil, errors.New("cannot connect to rabbit")
		}

		fmt.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}

	fmt.Println("Connected to RabbitMQ!")
	return rabbitConn, nil
}
