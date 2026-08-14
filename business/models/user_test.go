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

	u.Role = "user-admin"
	assert.Nil(t, u.ValidateRole())
	assert.False(t, u.IsAdmin())
	assert.True(t, u.IsUserManager())

	u.Role = ""
	assert.EqualError(t, u.ValidateRole(), ErrInvalidUserRole.Error())

	u.Role = "abc"
	assert.EqualError(t, u.ValidateRole(), ErrInvalidUserRole.Error())

	u.Vat = "N"
	assert.Nil(t, u.ValidateVat())

	u.Vat = "Y"
	assert.Nil(t, u.ValidateVat())

	u.Vat = ""
	assert.EqualError(t, u.ValidateVat(), ErrInvalidUserVat.Error())

	u.Vat = "abc"
	assert.EqualError(t, u.ValidateVat(), ErrInvalidUserVat.Error())

	u.CorporateName = "abc"
	u.Role = individual
	u.FirstName = "a"
	u.LastName = "b"
	assert.Equal(t, "a b", u.GetName())

	u.Role = corporate
	assert.Equal(t, "abc", u.GetName())
}

func TestUserClaims_Permissions(t *testing.T) {
	uc := &UserClaims{Role: "user-admin"}
	assert.False(t, uc.IsAdmin())
	assert.True(t, uc.IsUserManager())
	assert.False(t, uc.CanExportIncome())

	uc.Role = "admin"
	assert.True(t, uc.IsAdmin())
	assert.True(t, uc.IsUserManager())
	assert.True(t, uc.CanExportIncome())

	uc.Role = "individual"
	assert.False(t, uc.IsUserManager())
	assert.False(t, uc.CanExportIncome())
}

func TestUser_ApplyProfileUpdate(t *testing.T) {
	t.Run("does not mutate the source user", func(t *testing.T) {
		current := User{FirstName: "original"}
		updated := current
		updated.ApplyProfileUpdate(User{FirstName: "new"})
		assert.Equal(t, "original", current.FirstName)
		assert.Equal(t, "New", updated.FirstName)
	})

	t.Run("Bank account number with - and special char will create a bad batch file for bank system. This will fail the batch transfer process in the bank, causing the delay for all members to receive income. Therefore, we will remove - and special char from the bank account number!", func(t *testing.T) {
		current := User{BankAccountNumber: "000"}
		current.ApplyProfileUpdate(User{BankAccountNumber: "้1234-123-999‬"})
		assert.Equal(t, "1234123999", current.BankAccountNumber)
	})
}
