package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Saus")
	log.Println("extra saus")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("Error in connection")
	}

	defer conn.Close()

	fmt.Println("Succesfully connected to rabbit mq")

	chl, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer chl.Close()

	q, err := chl.QueueDeclare(
		"FlightQueue", // name of the queue
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		fmt.Println("Error in queue declaration")
		panic(err)
	}

	fmt.Println(q)

	err = chl.Publish(
		"",            // exchange
		"FlightQueue", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World"),
		},
	)

	if err != nil {
		fmt.Println("Error in publishing message")
		panic(err)
	}

	fmt.Println(" Pulished message to queue")
}
