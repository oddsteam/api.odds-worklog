package usecase

import "fmt"

func UpdateEvent(id int, firstName string, lastName string, dailyIncome int, workDate int, thaiCitizenId string, phone string, bank_no string, totalIncome float64, netIncome float64, netDailyIncome float64, wht float64, createdAt string, userId string, email string) string {
	return CreateEvent(id, firstName, lastName, dailyIncome, workDate, thaiCitizenId, phone, bank_no, totalIncome, netIncome, netDailyIncome, wht, createdAt, userId, email)
}

func CreateEvent(id int, firstName string, lastName string, dailyIncome int, workDate int, thaiCitizenId string, phone string, bank_no string, totalIncome float64, netIncome float64, netDailyIncome float64, wht float64, createdAt string, userId string, email string) string {
	format := `{
		"income":{
			"id":"%d",
			"totalIncome":"%f",
			"netIncome":"%f",
			"netDailyIncome":"%f",
			"workDate":%d,
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
