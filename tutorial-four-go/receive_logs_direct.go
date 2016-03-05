package main

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {

	tls := new(tls.Config)
	tls.InsecureSkipVerify = true

	// create connection
	conn, err := amqp.DialTLS("amqps://broker:CkY26kTuAyZT8r2@10.244.9.50:5671/", tls)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// create channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to connect to open a channel")

	err = ch.ExchangeDeclare(
		"logs_direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name,
			"logs_direct",
			s)

		err = ch.QueueBind(
			q.Name,
			s,
			"logs_direct",
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
