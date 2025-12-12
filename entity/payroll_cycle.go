package entity

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type PayrollCycle struct {
	records []*models.Income
	loans   models.StudentLoanList
}

func NewIncomes(records []*models.Income, loans models.StudentLoanList) *PayrollCycle {
	return &PayrollCycle{
		records: records,
		loans:   loans,
	}
}

func NewIncomesWithoutLoans(records []*models.Income) *PayrollCycle {
	return NewIncomes(records, models.StudentLoanList{})
}

func (ics *PayrollCycle) FindByUserID(id string) *models.Income {
	for _, e := range ics.records {
		if id == e.UserID {
			return e
		}
	}
	return &models.Income{}
}

func (ics *PayrollCycle) ProcessRecords(process func(index int, i Payroll) [][]string) ([][]string, []string) {
	strWrite := make([][]string, 0)
	updatedIncomeIds := []string{}
	for index, e := range ics.records {
		income := *e
		if income.ID.Hex() != "" {
			updatedIncomeIds = append(updatedIncomeIds, income.ID.Hex())
			loan := ics.loans.FindLoan(income.BankAccountName)
			i := NewIncomeFromRecord(income)
			i.SetLoan(&loan)
			rows := process(index, *i)
			strWrite = append(strWrite, rows...)
		}
	}
	return strWrite, updatedIncomeIds
}
