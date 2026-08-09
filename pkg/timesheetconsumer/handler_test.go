package timesheetconsumer

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

type fakeAcker struct {
	acked       bool
	ackMultiple bool
	nacked      bool
	nackRequeue bool
}

func (f *fakeAcker) Ack(multiple bool) error {
	f.acked = true
	f.ackMultiple = multiple
	return nil
}

func (f *fakeAcker) Nack(multiple, requeue bool) error {
	f.nacked = true
	f.nackRequeue = requeue
	return nil
}

func validEventBody() []byte {
	return []byte(`{
		"event_type": "timesheet.monthly_summary",
		"year": 2026, "month": 6,
		"summary_at": "2026-07-10T15:31:10+07:00",
		"employee": { "email": "employee@odds.team", "english_name": "Jane Doe" },
		"sites": []
	}`)
}

func TestHandleDelivery(t *testing.T) {
	t.Run("acks on successful sync", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		uc.EXPECT().SyncFromEvent(gomock.Any()).Return(nil)
		acker := &fakeAcker{}

		HandleDelivery(acker, validEventBody(), uc)

		assert.True(t, acker.acked)
		assert.False(t, acker.nacked)
	})

	t.Run("acks and drops on malformed JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		acker := &fakeAcker{}

		HandleDelivery(acker, []byte("not json"), uc)

		assert.True(t, acker.acked)
		assert.False(t, acker.nacked)
	})

	t.Run("acks and drops on unmatched user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		uc.EXPECT().SyncFromEvent(gomock.Any()).Return(usecases.ErrTimesheetUserNotFound)
		acker := &fakeAcker{}

		HandleDelivery(acker, validEventBody(), uc)

		assert.True(t, acker.acked)
		assert.False(t, acker.nacked)
	})

	t.Run("nacks and requeues on infra error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		uc.EXPECT().SyncFromEvent(gomock.Any()).Return(errors.New("mongo write failed"))
		acker := &fakeAcker{}

		HandleDelivery(acker, validEventBody(), uc)

		assert.False(t, acker.acked)
		assert.True(t, acker.nacked)
		assert.True(t, acker.nackRequeue)
	})

	t.Run("recovers a panic and acks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		uc.EXPECT().SyncFromEvent(gomock.Any()).DoAndReturn(func(models.TimesheetMonthlySummaryEvent) error {
			panic("boom")
		})
		acker := &fakeAcker{}

		assert.NotPanics(t, func() {
			HandleDelivery(acker, validEventBody(), uc)
		})
		assert.True(t, acker.acked)
		assert.False(t, acker.nacked)
	})
}
