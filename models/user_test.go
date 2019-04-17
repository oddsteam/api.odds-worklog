package models

import (
	"testing"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	u := new(User)
	u.FirstName = "Tester"
	u.LastName = "Super"
	assert.False(t, u.IsFullnameEmpty())

	u.Role = "admin"
	assert.True(t, u.IsAdmin())

	u.Role = ""
	assert.False(t, u.IsAdmin())

	u.Role = "individual"
	assert.False(t, u.IsAdmin())

	u.Role = "corporate"
	assert.False(t, u.IsAdmin())

	assert.Equal(t, "Tester Super", u.GetFullname())

	u.Email = "test@abc.com"
	assert.Equal(t, "test@abc.com", u.GetEmail())

	u.FirstName = ""
	assert.True(t, u.IsFullnameEmpty())

	u.FirstName = "Tester"
	u.LastName = ""
	assert.True(t, u.IsFullnameEmpty())

	u.Role = "admin"
	assert.Nil(t, u.ValidateRole())

	u.Role = "corporate"
	assert.Nil(t, u.ValidateRole())

	u.Role = "individual"
	assert.Nil(t, u.ValidateRole())

	u.Role = ""
	assert.EqualError(t, u.ValidateRole(), utils.ErrInvalidUserRole.Error())

	u.Role = "abc"
	assert.EqualError(t, u.ValidateRole(), utils.ErrInvalidUserRole.Error())

	u.Vat = "N"
	assert.Nil(t, u.ValidateVat())

	u.Vat = "Y"
	assert.Nil(t, u.ValidateVat())

	u.Vat = ""
	assert.EqualError(t, u.ValidateVat(), utils.ErrInvalidUserVat.Error())

	u.Vat = "abc"
	assert.EqualError(t, u.ValidateVat(), utils.ErrInvalidUserVat.Error())

	u.CorporateName = "abc"
	u.Role = individual
	u.FirstName = "a"
	u.LastName = "b"
	assert.Equal(t, "a b", u.GetName())

	u.Role = corporate
	assert.Equal(t, "abc", u.GetName())
}
