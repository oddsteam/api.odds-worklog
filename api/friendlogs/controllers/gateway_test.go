package controllers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/api/friendlogs/controllers"
)

func TestGatewayAddIncome(t *testing.T) {
	t.Run("panic should not cause subsriber to exit or the other events will be lost", func(t *testing.T) {
		assert.NotPanics(t, func() {
			controllers.CreateIncome(nil, "")
		})
	})
}

func TestGatewayUpdateIncome(t *testing.T) {
	t.Run("panic should not cause subsriber to exit or the other events will be lost", func(t *testing.T) {
		assert.NotPanics(t, func() {
			controllers.UpdateIncome(nil, "")
		})
	})
}
