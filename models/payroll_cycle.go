package models

type PayrollCycle struct {
	records []*Income
	loans   StudentLoanList
}

func NewPayrollCycle(records []*Income, loans StudentLoanList) *PayrollCycle {
	return &PayrollCycle{
		records: records,
		loans:   loans,
	}
}

func NewPayrollCycleWithoutLoans(records []*Income) *PayrollCycle {
	return NewPayrollCycle(records, StudentLoanList{})
}

func (pc *PayrollCycle) FindByUserID(id string) *Income {
	for _, e := range pc.records {
		if id == e.UserID {
			return e
		}
	}
	return &Income{}
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
