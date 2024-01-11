package consumer

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USER")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {
	consume()
}

func consume() {
	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open connection to RabbitMQ", err)
	}

	chl, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open the channel", err)
	}

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

	fmt.Println("Channel and queue established!")

	defer conn.Close()
	defer chl.Close()

	msgs, err := chl.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to register consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Recieved a message: %s", d.Body)
			d.Ack(false)
		}
	}()

	fmt.Println("Running...")

	<-forever
}
