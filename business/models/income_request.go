package models

type IncomeReq struct {
	Note          string `json:"note"`
	WorkDate      string `json:"workDate"`
	SpecialIncome string `json:"specialIncome"`
	WorkingHours  string `json:"workingHours"`
}



