package timesheetconsumer

import (
	"encoding/json"
	"errors"
	"log"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
)

// Acker is satisfied by github.com/rabbitmq/amqp091-go's Delivery.
type Acker interface {
	Ack(multiple bool) error
	Nack(multiple, requeue bool) error
}

func HandleDelivery(d Acker, body []byte, uc usecases.ForSyncingIncomeFromTimesheet) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("timesheetconsumer: panic recovered, dropping message: %v", r)
			d.Ack(false)
		}
	}()

	var evt models.TimesheetMonthlySummaryEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("timesheetconsumer: malformed event body, dropping: %v", err)
		d.Ack(false)
		return
	}

	err := uc.SyncFromEvent(evt)
	switch {
	case err == nil:
		d.Ack(false)
	case errors.Is(err, usecases.ErrTimesheetUserNotFound):
		log.Printf("timesheetconsumer: %v (email=%s, year=%d, month=%d), dropping", err, evt.Employee.Email, evt.Year, evt.Month)
		d.Ack(false)
	default:
		log.Printf("timesheetconsumer: sync failed (email=%s, year=%d, month=%d), requeueing: %v", evt.Employee.Email, evt.Year, evt.Month, err)
		d.Nack(false, true)
	}
}
