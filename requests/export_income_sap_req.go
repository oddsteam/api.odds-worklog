package requests

type ExportInComeSAPReq struct {
	Role          string `json:"role"`
	DateEffective string `json:"date_effective"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
}
