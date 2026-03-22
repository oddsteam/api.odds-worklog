package reminder_test

import (
	"os"
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
	origReceivers := os.Getenv("SMTP_REMINDER_RECEIVERS")
	defer os.Setenv("SMTP_REMINDER_RECEIVERS", origReceivers)
	os.Setenv("SMTP_REMINDER_RECEIVERS", "juacompe+worklog@odds.team,nalada@odds.team")

	u := createMockUser()
	m := reminder.CreateMailMessage(u, "id_copy_file_path.pdf")
	expected := "nalada@odds.team"

	if len(m.To) < 2 || m.To[1] != expected {
		t.Errorf("expected %v but got %v", expected, m.To)
	}
}

func TestMailIsAlsoSentToJuaToMonitor(t *testing.T) {
	origReceivers := os.Getenv("SMTP_REMINDER_RECEIVERS")
	defer os.Setenv("SMTP_REMINDER_RECEIVERS", origReceivers)
	os.Setenv("SMTP_REMINDER_RECEIVERS", "juacompe+worklog@odds.team,nalada@odds.team")

	u := createMockUser()
	m := reminder.CreateMailMessage(u, "id_copy_file_path.pdf")
	expected := "juacompe+worklog@odds.team"

	if len(m.To) < 1 || m.To[0] != expected {
		t.Errorf("expected %v but got %v", expected, m.To)
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

func TestNewSender_FailsWhenSMTPUserMissing(t *testing.T) {
	origUser := os.Getenv("SMTP_USER")
	origPass := os.Getenv("SMTP_PASSWORD")
	defer func() {
		os.Setenv("SMTP_USER", origUser)
		os.Setenv("SMTP_PASSWORD", origPass)
	}()

	os.Unsetenv("SMTP_USER")
	os.Setenv("SMTP_PASSWORD", "test")

	_, err := reminder.NewSender()
	if err == nil || !strings.Contains(err.Error(), "SMTP_USER") {
		t.Errorf("expected error about SMTP_USER, got: %v", err)
	}
}

func TestNewSender_FailsWhenSMTPPasswordMissing(t *testing.T) {
	origUser := os.Getenv("SMTP_USER")
	origPass := os.Getenv("SMTP_PASSWORD")
	defer func() {
		os.Setenv("SMTP_USER", origUser)
		os.Setenv("SMTP_PASSWORD", origPass)
	}()

	os.Setenv("SMTP_USER", "test@example.com")
	os.Unsetenv("SMTP_PASSWORD")

	_, err := reminder.NewSender()
	if err == nil || !strings.Contains(err.Error(), "SMTP_PASSWORD") {
		t.Errorf("expected error about SMTP_PASSWORD, got: %v", err)
	}
}

func TestNewSender_SucceedsWhenConfigPresent(t *testing.T) {
	origUser := os.Getenv("SMTP_USER")
	origPass := os.Getenv("SMTP_PASSWORD")
	origHost := os.Getenv("SMTP_HOST")
	origPort := os.Getenv("SMTP_PORT")
	defer func() {
		os.Setenv("SMTP_USER", origUser)
		os.Setenv("SMTP_PASSWORD", origPass)
		os.Setenv("SMTP_HOST", origHost)
		os.Setenv("SMTP_PORT", origPort)
	}()

	os.Setenv("SMTP_USER", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "testpassword")
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")

	sender, err := reminder.NewSender()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if sender == nil {
		t.Error("expected non-nil sender")
	}
}
