package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	broker "github.com/htaidirt/rabbit-go"
	"github.com/streadway/amqp"
)

func onFatalError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Setup the web router
	r := chi.NewRouter()

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

	rand.Seed(time.Now().UnixNano())
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		addTask := broker.AddTask{Number1: rand.Intn(999), Number2: rand.Intn(999)}
		body, err := json.Marshal(addTask)
		onFatalError(err, "Error encoding JSON")

		err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
		onFatalError(err, "Error publishing message")

		log.Printf("Published AddTask: %d + %d", addTask.Number1, addTask.Number2)
		w.Write([]byte("done"))
	})
	http.ListenAndServe(":3000", r)
}
