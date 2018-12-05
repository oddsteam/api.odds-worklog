package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	u := new(User)
	u.FirstName = "Tester"
	u.LastName = "Super"
	assert.False(t, u.IsFullnameEmpty())

	u.Email = "suphakrit@odds.team"
	assert.True(t, u.IsAdmin())

	u.Email = "jin@odds.team"
	assert.True(t, u.IsAdmin())

	u.Email = "roof@odds.team"
	assert.True(t, u.IsAdmin())

	u.Email = "a@odds.team"
	assert.False(t, u.IsAdmin())

	assert.Equal(t, "Tester Super", u.GetFullname())

	u.FirstName = ""
	assert.True(t, u.IsFullnameEmpty())

	u.FirstName = "Tester"
	u.LastName = ""
	assert.True(t, u.IsFullnameEmpty())
}
