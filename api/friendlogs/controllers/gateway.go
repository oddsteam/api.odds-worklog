package controllers

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/friendlogs/usecase"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func CreateIncome(s *mongo.Session, incomeCreatedEvent string) {
	record := usecase.NewUsecase().AddIncome(incomeCreatedEvent)
	err := income.NewRepository(s).AddIncome(&record)
	utils.FailOnError(err, "Fail to save event")
}

func UpdateIncome(s *mongo.Session, incomeUpdatedEvent string) {
	r := income.NewRepository(s)
	start, end := utils.GetStartDateAndEndDate(time.Now())
	incomes, err := r.GetAllIncomeByRoleStartDateAndEndDate("individual", start, end)
	utils.FailOnError(err, "Fail to save event")
	record := usecase.NewUsecase().UpdateIncome(incomes, incomeUpdatedEvent)
	err = income.NewRepository(s).UpdateIncome(&record)
	utils.FailOnError(err, "Fail to save event")
}
