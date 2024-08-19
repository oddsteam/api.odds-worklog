package main

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
	"gitlab.odds.team/worklog/api.odds-worklog/api/friendlogs/controllers"
	"gitlab.odds.team/worklog/api.odds-worklog/api/queue"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

func main() {
	conn := queue.Connect()
	defer conn.Close()

	session := mongo.Setup()
	defer session.Close()

	ch := queue.GetChannel(conn)
	defer ch.Close()

	createEvents := subscribe(ch, "incomes_created")

	forever := make(chan bool)

	go func() {
		for e := range createEvents {
			controllers.CreateIncome(session, string(e.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func subscribe(ch *amqp091.Channel, queueName string) <-chan amqp091.Delivery {
	q := queue.DeclareQueue(ch, "incomes_created", true)
	return queue.Subscribe(ch, q.Name)
}
