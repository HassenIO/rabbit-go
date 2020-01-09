package main

import (
	"encoding/json"
	"log"
	"os"

	broker "github.com/htaidirt/rabbit-go"
	"github.com/streadway/amqp"
)

func onFatalError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Setting up a RabbitMQ Connexion
	conn, err := amqp.Dial(broker.Config.AMQPConnectionURL)
	onFatalError(err, "Can't connect to AMQP")
	defer conn.Close()

	// Setting up a RabbitMQ Channel
	amqpChannel, err := conn.Channel()
	onFatalError(err, "Can't create channel")
	defer amqpChannel.Close()

	// Tell RabbitMQ which queue we are interested in.
	// QueueDeclare declares a queue to hold messages and deliver to consumers.
	// Declaring creates a queue if it doesn't already exist, or ensures that an existing queue matches the same parameters.
	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	onFatalError(err, "Can't create `add` queue")

	// Qos controls how many messages or how many bytes the server will try to keep on the network for consumers before receiving delivery acks.
	// The intent of Qos is to make sure the network buffers stay full between the server and client.
	err = amqpChannel.Qos(1, 0, false)
	onFatalError(err, "Can't configure QoS")

	// `Consume` method immediately starts delivering queued messages. Returns (<-chan Delivery, error)
	// Consumers must range over the returned chan to ensure all deliveries are received.
	// Unreceived deliveries will block all methods on the same connection.
	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	onFatalError(err, "Could not register consumer")

	// Create a channel to wait indefinitely
	stopChan := make(chan bool)
	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Consumer received a message: %s", d.Body)

			addTask := &broker.AddTask{}
			err := json.Unmarshal(d.Body, addTask)

			if err != nil {
				log.Printf("Error unmarshelling sent JSON: %s", err)
			}

			result := addTask.Number1 + addTask.Number2
			log.Println("Processing Calculation...")
			log.Printf("%d + %d = %d", addTask.Number1, addTask.Number2, result)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowleding: %s", err)
			} else {
				log.Printf("Finished consuming :)")
			}
		}
	}()
	<-stopChan // Never gets filled
}
