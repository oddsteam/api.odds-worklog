package main

import (
	"log"

	"gitlab.odds.team/worklog/api.odds-worklog/api/queue"
)

func main() {
	conn := queue.Connect()
	defer conn.Close()

	ch := queue.GetChannel(conn)
	defer ch.Close()

	q := queue.DeclareQueue(ch, "incomes_created", true)

	msgs := queue.Subscribe(ch, q.Name)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
