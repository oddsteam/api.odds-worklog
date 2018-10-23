package models

type IncomeInfo struct {
	User       User   `bson:"user" json:"user"`
	SubmitDate string `bson:"submitDate" json:"submitDate"`
	Status     string `bson:"status" json:"status"`
}
