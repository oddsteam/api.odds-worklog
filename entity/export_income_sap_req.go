package entity

type ExportInComeSAPReq struct {
	Role          string `json:"role"`
	DateEffective string `json:"dateEffective"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
}
