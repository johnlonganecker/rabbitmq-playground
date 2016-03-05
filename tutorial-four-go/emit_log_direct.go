package main

import (
	"crypto/tls"
	"log"
	"os"
	"strings"

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

	// create exchange
	ch.ExchangeDeclare(
		"logs_direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to connect to open a channel")

	body := bodyFrom(os.Args)
	// publish to channel
	err = ch.Publish(
		"logs_direct",
		severityFrom(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
