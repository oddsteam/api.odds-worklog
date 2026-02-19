package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

func TestNewUserDoesNotCauseSideEffectInTheGivenUserData(t *testing.T) {
	data := models.User{}
	data.FirstName = "original"
	u := NewUser(data)
	u.data.FirstName = "new"
	assert.Equal(t, "original", data.FirstName)
	assert.Equal(t, "new", u.data.FirstName)
}
