# Release note: income_from_timesheet

Summary of all commits since tag `v5.1.1` up to now (2026-09-06), split into what shipped in the previous round and what's in this round.

## Previous round (merged incrementally since v5.1.1)

- **Opened the timesheet inbox** — `GET /v1/timesheet-event-logs` shows historical events from the timesheet service (any logged-in user can view it, not admin-only like the SAP failure log)
- **Fixed inbox showing 0 working days/OT + blank site/employee names** — field tag mismatch (the JSON the web reads is camelCase, but the struct shared with the RabbitMQ consumer is snake_case from the publisher side). Fixed by mapping to a separate response type in the api layer instead of changing the tags on the struct the consumer uses
- **Added the income_from_timesheet system** — consumes the `timesheet.monthly_summary.published` event from RabbitMQ and upserts into the `income_from_timesheet` collection using the same payroll calculation as regular `income` (`CreatePayroll`/`UpdatePayroll`), kept in a strictly separate collection that doesn't touch existing `income`/`user`
- **Dual-write from manual add/edit income** — every time a user manually enters/edits income via the form, the system also mirrors the data to `income_from_timesheet` (`mirrorIncomeToTimesheet`) so this collection has complete data for everyone, not just people coming from timesheet
- **Export endpoints for income_from_timesheet** — added exports in multiple formats (CSV/SAP/PEAK) and refactored the export usecase to use a shared source adapter with the existing one
- **Notes can now be stored on both income and income_from_timesheet**
- **Moved PEAK Code to be stored on income itself** — no longer needs to join with `user` at export time
- **Stored site name on income itself** — same reason, no join needed
- **WHT rate 3% → 5%**
- **Refactor**: merged the SAP export failure ports (driving/driven, which looked identical) into a single interface

## This round (2026-09-06)

- **specialIncome no longer hardcoded to "0"** — when syncing from a timesheet event, the special hourly rate is now calculated from `user.DailyIncome / 8` (`business/usecases/sync_income_from_timesheet.go`)
- **Fixed OT unit conversion from timesheet** — `OvertimeDays` sent by timesheet is in **days**, now multiplied by 8 before storing into worklog's `WorkingHours` so the unit is **hours**, matching what payroll actually uses for calculation, and matching the manual add/edit side where users already enter hours (so it doesn't affect that side's dual-write)
- Updated the related unit tests to match the new behavior — `go build ./...` and `go test ./...` pass fully (422 tests)

## Pending

- Migrate data from `income_from_timesheet` to `income` so it shows up in history — waiting on Pi Jua
- Sync data to `income` for real — waiting on Pi Jua
- Remove the toggle and remove the dual-write (saving to `income` + `income_from_timesheet` simultaneously) — can only be done after data is fully synced

## Needs to be communicated to users

- If someone's OT rate differs from the normal rate (`dailyRate / 8`), the user needs to go fill in/adjust the hours themselves via the add/edit income form — auto-detect/adjust was not implemented in this round
