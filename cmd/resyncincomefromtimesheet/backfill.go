package main

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// matchTolerance is how far apart a row's insert time and an event's receivedAt may be and still
// be taken as the same delivery. The sync writes the event log and the row inside one message
// handler, so the real gap is milliseconds; the slack is there for a slow write, not for guessing.
const matchTolerance = time.Minute

// periodIndex answers which period a legacy row belongs to. Rows written before the fix carry no
// period at all, and their submitDate is the moment they were written rather than the month they
// report on, so the period has to be recovered from the event log the row was written from.
type periodIndex map[string][]*models.TimesheetEventLog

func newPeriodIndex(logs []*models.TimesheetEventLog) periodIndex {
	ix := make(periodIndex)
	for _, l := range logs {
		if l == nil {
			continue
		}
		ix[l.Employee.Email] = append(ix[l.Employee.Email], l)
	}
	return ix
}

// resolvePeriod returns the period to stamp on a legacy row. fromEvent reports whether it came from
// a matched timesheet delivery; when it is false the row is one the income form mirrored into this
// collection, and its own submitDate month is used — the period the old lookup already treated it
// as being in, so nothing moves.
func (ix periodIndex) resolvePeriod(row *models.IncomeFromTimesheet) (int, time.Month, bool) {
	insertedAt := row.ID.Timestamp()

	var best *models.TimesheetEventLog
	bestGap := matchTolerance
	for _, l := range ix[row.Email] {
		gap := insertedAt.Sub(l.ReceivedAt)
		if gap < 0 {
			gap = -gap
		}
		if gap <= bestGap {
			best, bestGap = l, gap
		}
	}

	if best != nil {
		return best.Year, time.Month(best.Month), true
	}

	submitted := row.SubmitDate.UTC()
	return submitted.Year(), submitted.Month(), false
}

// newestByLastUpdate picks the row to keep out of a set of duplicates for one period: the one
// written last, which carries whatever the most recent sync or manual edit left behind.
func newestByLastUpdate(rows []*models.IncomeFromTimesheet) *models.IncomeFromTimesheet {
	var keep *models.IncomeFromTimesheet
	for _, row := range rows {
		if keep == nil || row.LastUpdate.After(keep.LastUpdate) {
			keep = row
		}
	}
	return keep
}
