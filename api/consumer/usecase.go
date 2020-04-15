package consumer

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct {
	repo Repository
}

func NewUsecase(r Repository) Usecase {
	return &usecase{r}
}

func (u *usecase) GetByClientID(cid string) (*models.Consumer, error) {
	u.migrateClientID() // TODO: Remove this after migration
	consumer, err := u.repo.GetByClientID(cid)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// TODO: Remove this after migration
func (u *usecase) migrateClientID() {
	_, err := u.repo.GetByClientID("956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com")
	if err == utils.ErrInvalidConsumer {
		u.repo.Create(&models.Consumer{ClientID: "956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com"})
	}
}
