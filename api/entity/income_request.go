package entity

type IncomeReq struct {
	Note          string `bson:"note" json:"note"`
	WorkDate      string `bson:"workDate" json:"workDate"`
	SpecialIncome string `bson:"specialIncome" json:"specialIncome"`
	WorkingHours  string `bson:"workingHours" json:"workingHours"`
}
