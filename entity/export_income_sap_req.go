package entity

import (
	"errors"
	"time"
)

type ExportInComeSAPReq struct {
	Role          string `json:"role"`
	DateEffective string `json:"dateEffective"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
}

func (req *ExportInComeSAPReq) ParseDates() (time.Time, time.Time, time.Time, error) {
	startDate, err := req.getStartDate()
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}
	endDate, err := req.getEndDate()
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}
	dateEff, err := req.getDateEffective()
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}
	return startDate, endDate, dateEff, nil
}

func (req *ExportInComeSAPReq) getStartDate() (time.Time, error) {
	startDate, err := time.Parse("01/2006", req.StartDate)
	if err != nil {
		return time.Time{}, errors.New("startDate")
	}
	return startDate, nil
}

func (req *ExportInComeSAPReq) getEndDate() (time.Time, error) {
	endDate, err := time.Parse("01/2006", req.EndDate)
	if err != nil {
		return time.Time{}, errors.New("endDate")
	}
	endDate = endDate.AddDate(0, 1, 0)
	return endDate, nil
}

func (req *ExportInComeSAPReq) getDateEffective() (time.Time, error) {
	dateEff, err := time.Parse("02/01/2006", req.DateEffective)
	if err != nil {
		return time.Time{}, errors.New("dateEffective")
	}
	return dateEff, nil
}
