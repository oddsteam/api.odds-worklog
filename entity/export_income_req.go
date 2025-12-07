package entity

type ExportInComeReq struct {
	Role      string `json:"role"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
