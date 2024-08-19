package main

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
	"gitlab.odds.team/worklog/api.odds-worklog/api/friendlogs/usecase"
	"gitlab.odds.team/worklog/api.odds-worklog/api/queue"
)

func main() {
	conn := queue.Connect()
	defer conn.Close()

	ch := queue.GetChannel(conn)
	defer ch.Close()

	createdAt1 := "2024-07-26T06:26:25.531Z"
	updatedAt1 := "2024-07-27T06:26:25.531Z"
	createdAt2 := "2024-08-01T07:33:27.440Z"
	updatedAt2 := "2024-08-02T07:33:27.440Z"

	events := []string{
		usecase.CreateEvent(1, "Chi", "Sweethome", 750, 20,
			"123456789122", "+66912345678", "987654321",
			15375.0, 14913.75, 750.0, 461.25, createdAt1, createdAt1,
			"ba1357eb-20aa-4897-9759-658bf75e8429", "user1@example.com"),
		usecase.CreateEvent(2, "Yohei", "Yamada", 750, 10,
			"1234567890121", "0816543210", "0123456789",
			7500.0, 7275.0, 750.0, 225.0, createdAt2, createdAt2,
			"e82217a2-669a-4b0e-b98b-917e0ccfdf4c", "user2@example.com"),
	}

	publishMessages(ch, "incomes_created", events)

	events = []string{
		usecase.UpdateEvent(1, "Chi", "Sweethome", 750, 19,
			"123456789122", "+66912345678", "987654321",
			0.0, 0.0, 0.0, 0.0, createdAt1, updatedAt1,
			"ba1357eb-20aa-4897-9759-658bf75e8429", "user1@example.com"),
		usecase.UpdateEvent(2, "Yohei", "Yamada", 750, 10,
			"1234567890121", "0816543210", "0123456789",
			0.0, 0.0, 0.0, 0.0, createdAt2, updatedAt2,
			"e82217a2-669a-4b0e-b98b-917e0ccfdf4c", "user2@example.com"),
	}

	publishMessages(ch, "incomes_updated", events)
}

func publishMessages(ch *amqp091.Channel, queueName string, events []string) {
	q := queue.DeclareQueue(ch, queueName, true)

	for _, body := range events {
		queue.Publish(ch, q.Name, body)
		log.Printf(" [x] Sent %s", body)
	}
}
