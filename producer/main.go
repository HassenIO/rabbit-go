package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v7"
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

	// Create Redis Client to query database
	client := redis.NewClient(&redis.Options{
		Addr:     "database:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r.Get("/set/{name}/{age}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		age, err := strconv.Atoi(chi.URLParam(r, "age"))
		onFatalError(err, "Cannot convert age to integer")

		person := broker.PersonMsg{Name: name, Age: age}
		body, err := json.Marshal(person)
		onFatalError(err, "Error encoding JSON")

		err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
		onFatalError(err, "Error publishing message")

		log.Printf("Published PersonMsg: %s & %d", person.Name, person.Age)
		w.Write([]byte("done"))
	})

	r.Get("/get/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		log.Printf("URL parameter name=%s", name)

		val, err := client.Get(name).Result()
		if err != nil {
			log.Printf("Error getting data from database: %s", err)
			w.Write([]byte("KO"))
		} else {
			log.Printf("Got from database %s: %s", name, val)
			w.Write([]byte(val))
		}
	})

	http.ListenAndServe(":3000", r)
}
