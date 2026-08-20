package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// incomeFromTimesheetSource adapts the income_from_timesheet repository to the income source
// the export usecase reads from. The two collections carry the same Income payload, so
// unwrapping it here lets both sources share one export pipeline (CSV, PEAK and SAP alike)
// instead of duplicating it per source.
type incomeFromTimesheetSource struct {
	repo ForGettingIncomeFromTimesheetInTheMonth
}

func NewIncomeFromTimesheetSource(repo ForGettingIncomeFromTimesheetInTheMonth) ForGettingIncomeDataInTheMonth {
	return &incomeFromTimesheetSource{repo}
}

func (s *incomeFromTimesheetSource) GetAllIncomeByRoleStartDateAndEndDate(role string, startDate, endDate time.Time) ([]*models.Income, error) {
	records, err := s.repo.GetAllByRoleStartDateAndEndDate(role, startDate, endDate)
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
