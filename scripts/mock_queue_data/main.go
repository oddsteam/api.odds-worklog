package main

import (
	"log"

	"gitlab.odds.team/worklog/api.odds-worklog/api/friendlogs/usecase"
	"gitlab.odds.team/worklog/api.odds-worklog/api/queue"
)

func main() {
	conn := queue.Connect()
	defer conn.Close()

	ch := queue.GetChannel(conn)
	defer ch.Close()

	q := queue.DeclareQueue(ch, "incomes_created", true)

	events := []string{
		usecase.CreateEvent(1, "Chi", "Sweethome", 750, 20,
			"123456789122", "+66912345678", "987654321",
			15375.0, 14913.75, 750.0, 461.25, "2024-07-26T06:26:25.531Z",
			"ba1357eb-20aa-4897-9759-658bf75e8429", "user1@example.com"),
		usecase.CreateEvent(2, "Yohei", "Yamada", 750, 10,
			"1234567890121", "0816543210", "0123456789",
			7500.0, 7275.0, 750.0, 225.0, "2024-08-01T07:33:27.440Z",
			"e82217a2-669a-4b0e-b98b-917e0ccfdf4c", "user2@example.com"),
	}
	for _, body := range events {
		queue.Publish(ch, q.Name, body)
		log.Printf(" [x] Sent %s", body)
	}
}
