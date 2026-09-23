package file

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

func TestToSiteAllocationCSVWritesSiteAmountAndPercent(t *testing.T) {
	allocations := []models.SiteAllocation{
		{Site: "SCB", Amount: 1234567.5, Percent: 75},
		{Site: "KBANK", Amount: 411522.5, Percent: 25},
	}

	rows := ToSiteAllocationCSV(allocations)

	assert.Equal(t, []string{"SITE", "Amount", "Percent"}, rows[0])
	assert.Equal(t, []string{"SCB", "1,234,567.50", "75.00%"}, rows[1])
	assert.Equal(t, []string{"KBANK", "411,522.50", "25.00%"}, rows[2])
}

func TestToSiteAllocationCSVEndsWithATotalRowForReconciliation(t *testing.T) {
	allocations := []models.SiteAllocation{
		{Site: "SCB", Amount: 750, Percent: 75},
		{Site: "KBANK", Amount: 250, Percent: 25},
	}

	rows := ToSiteAllocationCSV(allocations)

	assert.Equal(t, []string{"TOTAL", "1,000.00", "100.00%"}, rows[len(rows)-1])
}

func TestSiteAllocationWriterRefusesToWriteAnEmptyReport(t *testing.T) {
	_, err := NewSiteAllocationCSVWriter().WriteFile("site_allocation", nil)

	assert.Error(t, err)
}

func TestSiteAllocationWriterWritesTheRowsToTheFile(t *testing.T) {
	allocations := []models.SiteAllocation{{Site: "SCB", Amount: 750, Percent: 100}}

	filename, err := NewSiteAllocationCSVWriter().WriteFile("site_allocation", allocations)
	assert.NoError(t, err)
	defer os.Remove(filename)

	f, err := os.Open(filename)
	assert.NoError(t, err)
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, [][]string{
		{"SITE", "Amount", "Percent"},
		{"SCB", "750.00", "100.00%"},
		{"TOTAL", "750.00", "100.00%"},
	}, rows)
}
