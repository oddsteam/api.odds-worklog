package controllers

import (
	"fmt"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/friendlogs/usecase"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func CreateIncome(s *mongo.Session, incomeCreatedEvent string) {
	defer handlePanic()
	saveIncome(s, incomeCreatedEvent, "Added")
}

func UpdateIncome(s *mongo.Session, incomeUpdatedEvent string) {
	defer handlePanic()
	saveIncome(s, incomeUpdatedEvent, "Updated")
}

func saveIncome(s *mongo.Session, event, action string) {
	er := NewRepository(s)
	er.Create(action, event)
	r := income.NewRepository(s)
	start, end := utils.GetStartDateAndEndDate(time.Now())
	incomes, err := r.GetAllIncomeByRoleStartDateAndEndDate("individual", start, end)
	utils.FailOnError(err, "Fail to retrieve incomes")
	record := usecase.NewUsecase().SaveIncome(incomes, event, action)
	if record.ID.Hex() == "" {
		err = r.AddIncomeOnSpecificTime(record, record.SubmitDate)
	} else {
		err = r.UpdateIncome(record)
	}
	utils.FailOnError(err, "Fail to save income")
}

func handlePanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
	}
}
