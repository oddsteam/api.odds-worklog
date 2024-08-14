package main

import (
	"fmt"
	"log"

	"gitlab.odds.team/worklog/api.odds-worklog/api/queue"
)

func main() {
	conn := queue.Connect()
	defer conn.Close()

	ch := queue.GetChannel(conn)
	defer ch.Close()

	q := queue.DeclareQueue(ch, "incomes_created", true)

	events := []string{
		createEvent(1, "Chi", "Sweethome", 750, 20,
			"123456789122", "+66912345678", "987654321",
			15375.0, 14913.75, 750.0, 461.25, "2024-07-26T06:26:25.531Z",
			"ba1357eb-20aa-4897-9759-658bf75e8429", "user1@example.com"),
		createEvent(2, "Yohei", "Yamada", 750, 10,
			"1234567890121", "0816543210", "0123456789",
			7500.0, 7275.0, 750.0, 225.0, "2024-08-01T07:33:27.440Z",
			"e82217a2-669a-4b0e-b98b-917e0ccfdf4c", "user2@example.com"),
	}
	for _, body := range events {
		queue.Publish(ch, q.Name, body)
		log.Printf(" [x] Sent %s", body)
	}
}

func createEvent(id int, firstName string, lastName string, dailyIncome int, workDate int, thaiCitizenId string, phone string, bank_no string, totalIncome float64, netIncome float64, netDailyIncome float64, wht float64, createdAt string, userId string, email string) string {
	format := `{
		"income":{
			"id":"%d",
			"totalIncome":"%f",
			"netIncome":"%f",
			"netDailyIncome":"%f",
			"workDate":"%d",
			"submitDate":null,
			"lastUpdate":null,
			"note":"",
			"vat":null,
			"wht":"%f",
			"specialIncome":null,
			"netSpecialIncome":null,
			"workingHours":null,
			"exportStatus":null,
			"created_at":"%s",
			"updated_at":"%s",
			"userId":"%s",
			"work_date":null
		},
		"registration":{
			"id":%d,
			"first_name":"%s",
			"last_name":"%s",
			"thai_citizen_id":"%s",
			"phone":"%s",
			"bank_no":"%s",
			"daily_income":"%d",
			"start_date":"2024-04-22T00:00:00.000Z",
			"created_at":"2024-07-26T06:26:18.565Z",
			"updated_at":"2024-07-26T06:26:18.565Z",
			"userId":"%s",
			"email":"%s"
		}
	}`
	return fmt.Sprintf(format, id, totalIncome, netIncome, netDailyIncome,
		workDate, wht, createdAt, createdAt, userId,
		id, firstName, lastName, thaiCitizenId, phone, bank_no,
		dailyIncome, userId, email)
}
