package main

import (
	"log"

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

	q := queue.DeclareQueue(ch, "incomes_created", true)

	msgs := queue.Subscribe(ch, q.Name)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			incomeCreatedEvent := string(d.Body)
			controllers.CreateIncome(session, incomeCreatedEvent)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
