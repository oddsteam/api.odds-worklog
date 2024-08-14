package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
)

func Connect() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	income.FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func GetChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	income.FailOnError(err, "Failed to open a channel")
	return ch
}

func DeclareQueue(ch *amqp.Channel, name string, durable bool) amqp.Queue {
	q, err := ch.QueueDeclare(
		name,
		durable,
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	income.FailOnError(err, "Failed to declare a queue")
	return q
}

func Publish(ch *amqp.Channel, routingKey string, body string) {
	err := ch.Publish(
		"",         // exchange
		routingKey, // routing key (queue name)
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	income.FailOnError(err, "Failed to publish a message")
}

func Subscribe(ch *amqp.Channel, routingKey string) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		routingKey, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}
	return msgs
}
