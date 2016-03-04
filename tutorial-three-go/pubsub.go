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

	conn, err := amqp.DialTLS("amqps://broker:CkY26kTuAyZT8r2@10.244.9.50:5671/", tls)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to connect to open a channel")

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	body := bodyFrom(os.Args)

	err = ch.Publish(
		"logs",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish")

	log.Printf(" [x] sent %s ", body)
}

func bodyFrom(args []string) string {
	var s string

	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
