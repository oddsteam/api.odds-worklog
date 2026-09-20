// Command resyncincomefromtimesheet repairs the income_from_timesheet collection after the period
// fix.
//
// Rows written before the fix carry no year/month, because the period only existed as an argument
// to the lookup and was never persisted. The lookup fell back to matching submitDate — the moment
// the sync ran, not the month being reported — so a re-delivered summary never found its row and
// was inserted again. The collection therefore holds one row per delivery instead of one row per
// user per period.
//
// This repairs the rows in place rather than rebuilding them. Nothing is recomputed: the payroll
// figures on every row are left exactly as they are, manual adjustments made through the income
// form included. Each legacy row gets the period it was really for — recovered from the event log
// entry it was written from — and duplicates of one period collapse to the row written last.
//
// It prints the plan and changes nothing unless -apply is given.
package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/joho/godotenv"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/repositories"
)

// periodSubmitDate mirrors the anchor the sync usecase writes, so a repaired row sits in the month
// it reports on for the period queries that read this collection.
func periodSubmitDate(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 12, 0, 0, 0, time.UTC)
}

type stamp struct {
	row       *models.IncomeFromTimesheet
	year      int
	month     time.Month
	fromEvent bool
}

func main() {
	apply := flag.Bool("apply", false, "write the changes; without it the command only reports what it would do")
	flag.Parse()

	_ = godotenv.Load()
	session := mongo.Setup()
	defer session.Close()

	rows, err := repositories.ListIncomeFromTimesheetRows(session)
	if err != nil {
		log.Fatalf("read income_from_timesheet: %v", err)
	}
	// List(0) is unlimited — the whole audit trail.
	logs, err := repositories.NewTimesheetEventLogLister(session).List(0)
	if err != nil {
		log.Fatalf("read timesheet_event_log: %v", err)
	}
	index := newPeriodIndex(logs)

	// Resolve a period for every row that has none, and group everything by period afterwards so
	// the duplicates show up.
	var stamps []stamp
	var fromEvent, fromSubmitDate int
	groups := make(map[string][]*models.IncomeFromTimesheet)
	for _, row := range rows {
		year, month := row.Year, time.Month(row.Month)
		if row.Year == 0 {
			var matched bool
			year, month, matched = index.resolvePeriod(row)
			stamps = append(stamps, stamp{row, year, month, matched})
			if matched {
				fromEvent++
			} else {
				fromSubmitDate++
			}
		}
		key := fmt.Sprintf("%s|%04d-%02d", row.UserID, year, month)
		groups[key] = append(groups[key], row)
	}

	var drop []*models.IncomeFromTimesheet
	var dupKeys []string
	for key, g := range groups {
		if len(g) < 2 {
			continue
		}
		dupKeys = append(dupKeys, key)
		keep := newestByLastUpdate(g)
		for _, row := range g {
			if row.ID != keep.ID {
				drop = append(drop, row)
			}
		}
	}
	sort.Strings(dupKeys)

	log.Printf("income_from_timesheet: %d rows, %d without a stored period", len(rows), len(stamps))
	log.Printf("period recovered from the event log: %d rows; read off submitDate (mirrored from the income form): %d rows", fromEvent, fromSubmitDate)
	log.Printf("periods holding duplicates: %d — %d rows to drop, keeping the most recently written one of each", len(dupKeys), len(drop))
	for i, key := range dupKeys {
		if i == 20 {
			log.Printf("  … and %d more", len(dupKeys)-20)
			break
		}
		log.Printf("  %s: %d rows", key, len(groups[key]))
	}

	if !*apply {
		log.Print("dry run: nothing written. re-run with -apply")
		return
	}

	for _, s := range stamps {
		if err := repositories.SetIncomeFromTimesheetPeriod(session, s.row.ID, s.year, s.month, periodSubmitDate(s.year, s.month)); err != nil {
			log.Fatalf("stamp period on %s: %v", s.row.ID.Hex(), err)
		}
	}
	log.Printf("stamped the period on %d rows", len(stamps))

	for _, row := range drop {
		if err := repositories.DeleteIncomeFromTimesheetRow(session, row.ID); err != nil {
			log.Fatalf("delete duplicate %s: %v", row.ID.Hex(), err)
		}
	}
	log.Printf("dropped %d duplicate rows", len(drop))

	if err := repositories.EnsureIncomeFromTimesheetPeriodIndex(session); err != nil {
		log.Fatalf("create unique period index: %v", err)
	}
	log.Print("unique index on userId+year+month is in place")
}
