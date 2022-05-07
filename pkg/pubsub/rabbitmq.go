package pubsub

import (
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	connectionString string
}

func NewRabbitMQ(connectionString string) *RabbitMQ {
	return &RabbitMQ{
		connectionString: connectionString,
	}
}

func (mq *RabbitMQ) Publish(topic, message string) error {
	log.Printf("publishing message \"%s\" to topic \"%s\"", message, topic)
	conn, err := amqp.Dial(mq.connectionString)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Println("Failed to declare a queue")
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Println("Failed to publish a message")
		return err
	}
	return nil
}

func (mq *RabbitMQ) Subscribe(topic string, callback func(string) error) error {
	connectionString := "amqps://gtscuqup:j9rcANXywRogZgAu8I7T7oE-3XueLvdq@jackal.rmq.cloudamqp.com/gtscuqup"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Println("Failed to declare a queue")
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println("Failed to register a consumer")
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			c_err := callback(string(d.Body))
			if c_err != nil {
				log.Println("Failed to process message")
				continue
			}
			d.Ack(false)
			log.Printf("Done")

		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}
