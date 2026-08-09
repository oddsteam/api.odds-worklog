package main

import (
	"context"
	"flag"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/config"
)

func main() {
	file := flag.String("file", "", "path to a JSON file containing a timesheet.monthly_summary event")
	flag.Parse()
	if *file == "" {
		log.Fatal("missing required -file flag")
	}

	body, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalf("read %s: %v", *file, err)
	}

	c := config.Config()
	conn, err := amqp.Dial(c.RabbitMQURL)
	if err != nil {
		log.Fatalf("dial %s: %v", c.RabbitMQURL, err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("open channel: %v", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(c.RabbitMQExchange, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("declare exchange %s: %v", c.RabbitMQExchange, err)
	}

	err = ch.PublishWithContext(context.Background(), c.RabbitMQExchange, c.RabbitMQRoutingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		log.Fatalf("publish: %v", err)
	}

	log.Printf("published %s to exchange=%s routing_key=%s", *file, c.RabbitMQExchange, c.RabbitMQRoutingKey)
}
