package consumer

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	repo Repository
}

func NewUsecase(r Repository) Usecase {
	return &usecase{r}
}

func (u *usecase) GetByClientID(cid string) (*models.Consumer, error) {
	consumer, err := u.repo.GetByClientID(cid)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
