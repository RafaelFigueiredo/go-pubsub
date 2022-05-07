package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/rafaelfigueiredo/rabbitmq/pkg/pubsub"
)

var amqp_url = os.Getenv("AMQP_URL")

func main() {
	ps := pubsub.NewRabbitMQ(amqp_url)
	err := ps.Publish("deploy", "testapp")
	failOnError(err)

	ps.Subscribe("deploy", func(msg string) error {
		if msg == "tofail" {
			return errors.New("test error")
		}
		log.Println("doing some process")
		time.Sleep(1 * time.Second)
		return nil
	})
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
