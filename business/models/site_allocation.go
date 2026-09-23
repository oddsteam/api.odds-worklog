package models

import "sort"

// UnassignedSite catches income that no site can be derived for — the timesheet event
// carried no sites, or recorded no days on the ones it did. Such income is still reported
// so the breakdown adds up to the same total as the payroll export for the month.
const UnassignedSite = "(ไม่ระบุ site)"

// SiteAllocation is one row of the per-site income breakdown: how much of the month's
// income the site accounts for, and what share of the whole that is.
type SiteAllocation struct {
	Site    string
	Amount  float64
	Percent float64
}

// BuildSiteAllocations spreads each person's income across the sites they worked at.
// An Income carries only the person's monthly totals, so the per-site figures are derived
// from the day counts the timesheet event recorded per site.
func BuildSiteAllocations(records []*IncomeFromTimesheet) []SiteAllocation {
	amountBySite := map[string]float64{}

	for _, record := range records {
		dailyIncome, _ := StringToFloat64(record.DailyIncomeBeforeTax)
		specialIncome, _ := StringToFloat64(record.SpecialIncomeBeforeTax)
		if dailyIncome+specialIncome == 0 {
			continue
		}

		var totalWorkingDays, totalOvertimeDays float64
		for _, site := range record.Sites {
			totalWorkingDays += site.WorkingDays
			totalOvertimeDays += site.OvertimeDays
		}

		// Nothing to spread the income over: park it rather than divide by zero.
		if totalWorkingDays == 0 && totalOvertimeDays == 0 {
			amountBySite[UnassignedSite] += dailyIncome + specialIncome
			continue
		}

		// Main income follows the work days and special income the overtime days: the two
		// are earned at different rates, so one blended ratio would misattribute the money
		// of anyone whose OT is concentrated on a single site.
		for _, site := range record.Sites {
			amountBySite[site.ClientSite] += share(dailyIncome, site.WorkingDays, totalWorkingDays) +
				share(specialIncome, site.OvertimeDays, totalOvertimeDays)
		}
	}

	var total float64
	for _, amount := range amountBySite {
		total += amount
	}

	allocations := make([]SiteAllocation, 0, len(amountBySite))
	for site, amount := range amountBySite {
		allocations = append(allocations, SiteAllocation{
			Site:    site,
			Amount:  amount,
			Percent: amount / total * 100,
		})
	}
	// Biggest spend first, with the catch-all bucket pinned to the bottom and equal
	// amounts broken by name so the same month always exports the same file.
	sort.Slice(allocations, func(i, j int) bool {
		a, b := allocations[i], allocations[j]
		if (a.Site == UnassignedSite) != (b.Site == UnassignedSite) {
			return b.Site == UnassignedSite
		}
		if a.Amount != b.Amount {
			return a.Amount > b.Amount
		}
		return a.Site < b.Site
	})

	return allocations
}

// share is amount * part / whole, guarding the case where a record carries income but no
// days to spread it over.
func share(amount, part, whole float64) float64 {
	if whole == 0 {
		return 0
	}
	return amount * part / whole
}
