package usecase

import "fmt"

func UpdateEvent(id int, firstName, lastName string, dailyIncome int, workDate float64, thaiCitizenId, phone, bank_no string, totalIncome, netIncome, netDailyIncome, wht float64, createdAt, updatedAt, userId, email string) string {
	return CreateEvent(id, firstName, lastName, dailyIncome, workDate, thaiCitizenId, phone, bank_no, totalIncome, netIncome, netDailyIncome, wht, createdAt, updatedAt, userId, email)
}

func CreateEvent(id int, firstName, lastName string, dailyIncome int, workDate float64, thaiCitizenId, phone, bank_no string, totalIncome, netIncome, netDailyIncome, wht float64, createdAt, updatedAt, userId, email string) string {
	format := `{
		"income":{
			"id":"%d",
			"totalIncome":"%f",
			"netIncome":"%f",
			"netDailyIncome":"%f",
			"workDate": %f,
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
		workDate, wht, createdAt, updatedAt, userId,
		id, firstName, lastName, thaiCitizenId, phone, bank_no,
		dailyIncome, userId, email)
}
