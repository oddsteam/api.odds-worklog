package user

import (
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type User struct {
	data      *models.User
	DailyRate float64
}

func NewUser(data models.User) *User {
	return &User{data: &data, DailyRate: 0}
}

func (u *User) Parse() error {
	var err error
	u.DailyRate, err = utils.StringToFloat64(u.data.DailyIncome)
	if err != nil {
		return err
	}
	return nil
}
func (u *User) Role() string {
	return u.data.Role
}

func (u *User) IsVATRegistered() bool {
	return u.data.Vat == "Y"
}

func (u *User) prepareDataForUpdateFrom(m models.User) error {
	if m.FirstName != "" {
		u.data.FirstName = utils.ToFirstUpper(m.FirstName)
	}
	if m.LastName != "" {
		u.data.LastName = utils.ToFirstUpper(m.LastName)
	}
	if m.CorporateName != "" {
		u.data.CorporateName = m.CorporateName
	}
	if m.BankAccountName != "" {
		u.data.BankAccountName = m.BankAccountName
	}
	if m.BankAccountNumber != "" {
		u.data.BankAccountNumber = extractNumbers(m.BankAccountNumber)
	}
	if m.ThaiCitizenID != "" {
		u.data.ThaiCitizenID = m.ThaiCitizenID
	}
	if m.SlackAccount != "" {
		if err := utils.ValidateEmail(m.SlackAccount); err != nil {
			return errors.New("Invalid Slack acount")
		}
		u.data.SlackAccount = m.SlackAccount
	}
	if m.SiteID != "" {
		u.data.SiteID = m.SiteID
	}
	if m.Project != "" {
		u.data.Project = m.Project
	}
	if m.DailyIncome != "" {
		u.data.DailyIncome = m.DailyIncome
	}
	if m.Address != "" {
		u.data.Address = m.Address
	}
	if m.StartDate != "" {
		u.data.StartDate = m.StartDate
	}

	u.data.StatusTavi = m.StatusTavi
	u.data.Role = m.Role
	u.data.Vat = m.Vat
	u.data.Phone = m.Phone

	return nil
}
