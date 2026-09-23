package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSiteAllocationsSplitsDailyIncomeByWorkingDayShare(t *testing.T) {
	records := []*IncomeFromTimesheet{
		{
			Income: Income{DailyIncomeBeforeTax: "1000.00", SpecialIncomeBeforeTax: "0.00"},
			Sites: []SiteWork{
				{ClientSite: "SCB", WorkingDays: 15},
				{ClientSite: "KBANK", WorkingDays: 5},
			},
		},
	}

	allocations := BuildSiteAllocations(records)

	assert.Equal(t, []SiteAllocation{
		{Site: "SCB", Amount: 750, Percent: 75},
		{Site: "KBANK", Amount: 250, Percent: 25},
	}, allocations)
}

func TestBuildSiteAllocationsSplitsSpecialIncomeByOvertimeDayShare(t *testing.T) {
	// Work days and OT days sit on different sites, so a single blended ratio would
	// misattribute the OT money: it has to follow the overtime days alone.
	records := []*IncomeFromTimesheet{
		{
			Income: Income{DailyIncomeBeforeTax: "1000.00", SpecialIncomeBeforeTax: "400.00"},
			Sites: []SiteWork{
				{ClientSite: "SCB", WorkingDays: 20, OvertimeDays: 0},
				{ClientSite: "KBANK", WorkingDays: 0, OvertimeDays: 2},
			},
		},
	}

	allocations := BuildSiteAllocations(records)

	assert.Len(t, allocations, 2)
	assert.Equal(t, "SCB", allocations[0].Site)
	assert.InDelta(t, 1000, allocations[0].Amount, 0.001)
	assert.InDelta(t, 71.43, allocations[0].Percent, 0.01)
	assert.Equal(t, "KBANK", allocations[1].Site)
	assert.InDelta(t, 400, allocations[1].Amount, 0.001)
	assert.InDelta(t, 28.57, allocations[1].Percent, 0.01)
}

func TestBuildSiteAllocationsSumsTheSameSiteAcrossPeople(t *testing.T) {
	records := []*IncomeFromTimesheet{
		{
			Income: Income{DailyIncomeBeforeTax: "1000.00"},
			Sites:  []SiteWork{{ClientSite: "SCB", WorkingDays: 20}},
		},
		{
			Income: Income{DailyIncomeBeforeTax: "3000.00"},
			Sites:  []SiteWork{{ClientSite: "SCB", WorkingDays: 10}},
		},
	}

	allocations := BuildSiteAllocations(records)

	assert.Equal(t, []SiteAllocation{{Site: "SCB", Amount: 4000, Percent: 100}}, allocations)
}

func TestBuildSiteAllocationsPutsIncomeWithoutSitesInTheUnassignedBucket(t *testing.T) {
	// The money still left the company, so it has to show up somewhere — dropping it
	// would make the report disagree with the payroll export for the same month.
	records := []*IncomeFromTimesheet{
		{
			Income: Income{DailyIncomeBeforeTax: "600.00", SpecialIncomeBeforeTax: "400.00"},
			Sites:  nil,
		},
		{
			Income: Income{DailyIncomeBeforeTax: "1000.00"},
			Sites:  []SiteWork{{ClientSite: "SCB", WorkingDays: 20}},
		},
	}

	allocations := BuildSiteAllocations(records)

	assert.Equal(t, []SiteAllocation{
		{Site: "SCB", Amount: 1000, Percent: 50},
		{Site: UnassignedSite, Amount: 1000, Percent: 50},
	}, allocations)
}

func TestBuildSiteAllocationsKeepsIncomeWhoseSitesRecordNoDays(t *testing.T) {
	records := []*IncomeFromTimesheet{
		{
			Income: Income{DailyIncomeBeforeTax: "500.00"},
			Sites:  []SiteWork{{ClientSite: "SCB", WorkingDays: 0, OvertimeDays: 0}},
		},
	}

	allocations := BuildSiteAllocations(records)

	assert.Equal(t, []SiteAllocation{{Site: UnassignedSite, Amount: 500, Percent: 100}}, allocations)
}

func TestBuildSiteAllocationsReturnsNoRowsWhenThereIsNoIncome(t *testing.T) {
	assert.Empty(t, BuildSiteAllocations(nil))
}

func TestBuildSiteAllocationsBreaksEqualAmountsByNameSoExportsAreStable(t *testing.T) {
	records := []*IncomeFromTimesheet{
		{
			Income: Income{DailyIncomeBeforeTax: "1000.00"},
			Sites: []SiteWork{
				{ClientSite: "SCB", WorkingDays: 10},
				{ClientSite: "KBANK", WorkingDays: 10},
			},
		},
	}

	allocations := BuildSiteAllocations(records)

	assert.Equal(t, []string{"KBANK", "SCB"}, []string{allocations[0].Site, allocations[1].Site})
}

func TestBuildSiteAllocationsIgnoresRecordsWithNoIncomeToSplit(t *testing.T) {
	// A record with no money on it would otherwise contribute a zero row — and if every
	// record is like that, a division by a zero total turns every Percent into NaN.
	records := []*IncomeFromTimesheet{
		{
			Income: Income{DailyIncomeBeforeTax: "", SpecialIncomeBeforeTax: ""},
			Sites:  []SiteWork{{ClientSite: "SCB", WorkingDays: 20}},
		},
		{
			Income: Income{DailyIncomeBeforeTax: "0.00"},
			Sites:  nil,
		},
	}

	assert.Empty(t, BuildSiteAllocations(records))
}
