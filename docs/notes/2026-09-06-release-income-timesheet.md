# Release note: income_from_timesheet

Summary of all commits since tag `v5.1.1` up to now (2026-09-06), split into what shipped in the previous round and what's in this round.

## Previous round (merged incrementally since v5.1.1)

- **Opened the timesheet inbox** — `GET /v1/timesheet-event-logs` shows historical events from the timesheet service (any logged-in user can view it, not admin-only like the SAP failure log)
- **Fixed inbox showing 0 working days/OT + blank site/employee names** — field tag mismatch (the JSON the web reads is camelCase, but the struct shared with the RabbitMQ consumer is snake_case from the publisher side). Fixed by mapping to a separate response type in the api layer instead of changing the tags on the struct the consumer uses
- **Added the income_from_timesheet system** — consumes the `timesheet.monthly_summary.published` event from RabbitMQ and upserts into the `income_from_timesheet` collection using the same payroll calculation as regular `income` (`CreatePayroll`/`UpdatePayroll`), kept in a strictly separate collection that doesn't touch existing `income`/`user`

## This round (2026-09-06)

- **specialIncome no longer hardcoded to "0"** — when syncing from a timesheet event, the special hourly rate is now calculated from `user.DailyIncome / 8` (`business/usecases/sync_income_from_timesheet.go`)
- **Fixed OT unit conversion from timesheet** — `OvertimeDays` sent by timesheet is in **days**, now multiplied by 8 before storing into worklog's `WorkingHours` so the unit is **hours**, matching what payroll uses for calculation
- **Missing/invalid daily rate is now logged to error-logs** — if a user has OT hours from timesheet but no valid `DailyIncome`, special income can't be calculated; this now gets recorded via the SAP export failure log (`/error-logs`) instead of silently dropping the income

## Pending

- Migrate data from `income_from_timesheet` to `income` so it shows up in history — waiting on Pi Jua
- Sync data to `income` for real — waiting on Pi Jua
- Remove the toggle and remove the dual-write (saving to `income` + `income_from_timesheet` simultaneously) — can only be done after data is fully synced

## Needs to be communicated to users

- If someone's OT rate differs from the normal rate (`dailyRate / 8`), the user needs to go fill in/adjust the hours themselves via the add/edit income form — auto-detect/adjust was not implemented in this round
