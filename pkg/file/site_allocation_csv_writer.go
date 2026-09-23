package file

import (
	"encoding/csv"
	"errors"
	"fmt"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type siteAllocationCSVWriter struct{}

func NewSiteAllocationCSVWriter() *siteAllocationCSVWriter {
	return &siteAllocationCSVWriter{}
}

func (w *siteAllocationCSVWriter) WriteFile(name string, allocations []models.SiteAllocation) (string, error) {
	if len(allocations) == 0 {
		return "", errors.New("no data for export to CSV file")
	}

	file, filename, err := CreateFile(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(ToSiteAllocationCSV(allocations))
	csvWriter.Flush()
	return filename, nil
}

// ToSiteAllocationCSV renders the breakdown, closing with a TOTAL row so the report can be
// reconciled against the payroll export for the same period at a glance.
func ToSiteAllocationCSV(allocations []models.SiteAllocation) [][]string {
	rows := [][]string{{"SITE", "Amount", "Percent"}}

	var totalAmount, totalPercent float64
	for _, a := range allocations {
		totalAmount += a.Amount
		totalPercent += a.Percent
		rows = append(rows, []string{a.Site, siteAllocationAmount(a.Amount), siteAllocationPercent(a.Percent)})
	}

	return append(rows, []string{"TOTAL", siteAllocationAmount(totalAmount), siteAllocationPercent(totalPercent)})
}

func siteAllocationAmount(amount float64) string {
	return models.FormatCommas(models.FloatToString(amount))
}

func siteAllocationPercent(percent float64) string {
	return fmt.Sprintf("%.2f%%", percent)
}
