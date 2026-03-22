package reminder_test

import (
	"strings"
	"testing"

	"gitlab.odds.team/worklog/api.odds-worklog/api/reminder"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

func TestMailMessageShouldContainsBankAccountNameSoWeCanFindTheOldMailInMailboxUsingUserThaiName(t *testing.T) {
	u := createMockUser()
	m := reminder.CreateMailMessage(u, "id_copy_file_path.pdf")

	if !strings.Contains(m.Body, u.BankAccountName) {
		t.Errorf("Should contains %v but not (%v)", u.BankAccountName, m.Body)
	}
}

func TestMailIsSentToFinance(t *testing.T) {
	u := createMockUser()
	m := reminder.CreateMailMessage(u, "id_copy_file_path.pdf")
	expected := "nalada@odds.team"

	if m.To[1] != expected {
		t.Errorf("expected %v but got %v", expected, m.To[1])
	}
}

func TestMailIsAlsoSentToJuaToMonitor(t *testing.T) {
	u := createMockUser()
	m := reminder.CreateMailMessage(u, "id_copy_file_path.pdf")
	expected := "juacompe+worklog@odds.team"

	if m.To[0] != expected {
		t.Errorf("expected %v but got %v", expected, m.To[0])
	}
}

func createMockUser() models.User {
	return models.User{
		FirstName:       "FirstName",
		LastName:        "LastName",
		BankAccountName: "นาย ชื่อไทย นามสกุล",
		Email:           "mail@odds.team",
	}
}
