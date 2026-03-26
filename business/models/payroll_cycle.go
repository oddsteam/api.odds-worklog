package models

// IncomeRowMeta identifies the income record that produced one or more export rows (e.g. two SAP lines per income).
type IncomeRowMeta struct {
	IncomeID        string
	UserID          string
	BankAccountName string
}

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

func (pc *PayrollCycle) ProcessRecords(process func(index int, i Payroll) [][]string) ([][]string, []IncomeRowMeta) {
	strWrite := make([][]string, 0)
	meta := make([]IncomeRowMeta, 0)
	for index, e := range pc.records {
		income := *e
		if income.ID.Hex() != "" {
			meta = append(meta, IncomeRowMeta{
				IncomeID:        income.ID.Hex(),
				UserID:          income.UserID,
				BankAccountName: income.BankAccountName,
			})
			loan := pc.loans.FindLoan(income.BankAccountName)
			i := NewPayrollFromIncome(income)
			i.SetLoan(&loan)
			rows := process(index, *i)
			strWrite = append(strWrite, rows...)
		}
	}
	return strWrite, meta
}
