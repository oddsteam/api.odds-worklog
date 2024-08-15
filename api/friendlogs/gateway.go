package friendslog

import (
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func CreateIncome(s *mongo.Session, incomeCreatedEvent string) {
	record := NewUsecase().AddIncome(incomeCreatedEvent)
	err := income.NewRepository(s).AddIncome(&record)
	utils.FailOnError(err, "Fail to save event")
}
