package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// incomeFromTimesheetUserSource adapts the income_from_timesheet repository to
// ForReadingUserIncome, so the individual income list and status endpoints can read from
// either source through the same usecases, mirroring incomeFromTimesheetSource for export.
type incomeFromTimesheetUserSource struct {
	repo ForReadingIncomeFromTimesheetByUser
}

func NewIncomeFromTimesheetUserSource(repo ForReadingIncomeFromTimesheetByUser) ForReadingUserIncome {
	return &incomeFromTimesheetUserSource{repo}
}

func (s *incomeFromTimesheetUserSource) GetIncomeUserByYearMonth(id string, fromYear int, fromMonth time.Month) (*models.Income, error) {
	record, err := s.repo.GetByUserYearMonth(id, fromYear, fromMonth)
	if err != nil {
		return nil, err
	}
	return &record.Income, nil
}

func (s *incomeFromTimesheetUserSource) GetIncomeByUserIdAllMonth(userId string) ([]*models.Income, error) {
	records, err := s.repo.GetByUserIdAllMonth(userId)
	if err != nil {
		return nil, err
	}

	incomes := make([]*models.Income, 0, len(records))
	for _, record := range records {
		if record == nil {
			continue
		}
		incomes = append(incomes, &record.Income)
	}
	return incomes, nil
}
