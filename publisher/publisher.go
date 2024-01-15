package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
	"github.com/streadway/amqp"
	httpSwagger "github.com/swaggo/http-swagger"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USER")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {
	router := httprouter.New()

	var swaggerDoc = "http://localhost:5000/swagger/doc.json"

	// Reformat httpSwagger Handler function to expected format for net/http router
	swaggerHandler := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		httpSwagger.Handler(
			httpSwagger.URL(swaggerDoc), //The url pointing to API definition
		)
	}

	router.GET("/swagger/*", swaggerHandler)

	router.POST("/api/publish/:message", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		submit(w, r, p)
	})

	fmt.Println("Running")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func submit(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	message := p.ByName("message")

	fmt.Println("Recieved: " + message)

	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open connection to RabbitMQ", err)
	}

	defer conn.Close()

	chl, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open the channel", err)
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
		log.Fatalf("%s: %s", "Failed to open the queue", err)
	}

	err = chl.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish message to the queue", err)
	}

	fmt.Println("Post succesful!")
}
