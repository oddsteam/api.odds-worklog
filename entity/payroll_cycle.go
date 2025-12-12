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

func (pc *PayrollCycle) FindByUserID(id string) *models.Income {
	for _, e := range pc.records {
		if id == e.UserID {
			return e
		}
	}
	return &models.Income{}
}

func (pc *PayrollCycle) ProcessRecords(process func(index int, i Payroll) [][]string) ([][]string, []string) {
	strWrite := make([][]string, 0)
	updatedIncomeIds := []string{}
	for index, e := range pc.records {
		income := *e
		if income.ID.Hex() != "" {
			updatedIncomeIds = append(updatedIncomeIds, income.ID.Hex())
			loan := pc.loans.FindLoan(income.BankAccountName)
			i := NewPayrollFromIncome(income)
			i.SetLoan(&loan)
			rows := process(index, *i)
			strWrite = append(strWrite, rows...)
		}
	}
	return strWrite, updatedIncomeIds
}
