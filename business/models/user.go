package models

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Role              string        `bson:"role" json:"role"`
	FirstName         string        `bson:"firstName" json:"firstName"`
	LastName          string        `bson:"lastName" json:"lastName"`
	CorporateName     string        `bson:"corporateName" json:"corporateName,omitempty"`
	Email             string        `bson:"email" json:"email"`
	PeakCode          string        `bson:"peakCode" json:"peakCode"`
	BankAccountName   string        `bson:"bankAccountName" json:"bankAccountName"`
	BankAccountNumber string        `bson:"bankAccountNumber" json:"bankAccountNumber"`
	ThaiCitizenID     string        `bson:"thaiCitizenId" json:"thaiCitizenId,omitempty"`
	Vat               string        `bson:"vat" json:"vat,omitempty"`
	Transcript        string        `bson:"transcript" json:"transcript,omitempty"`
	SiteID            string        `bson:"siteId" json:"siteId,omitempty"`
	Project           string        `bson:"project" json:"project,omitempty"`
	ImageProfile      string        `bson:"imageProfile" json:"imageProfile,omitempty"`
	DegreeCertificate string        `bson:"degreeCertificate" json:"degreeCertificate,omitempty"`
	IDCard            string        `bson:"idCard" json:"idCard,omitempty"`
	Site              *Site         `bson:"-" json:"site,omitempty"`
	Create            time.Time     `bson:"create" json:"create"`
	LastUpdate        time.Time     `bson:"lastUpdate" json:"lastUpdate"`
	DailyIncome       string        `bson:"dailyIncome" json:"dailyIncome,omitempty"`
	Address           string        `bson:"address" json:"address,omitempty"`
	StatusTavi        bool          `bson:"statusTavi" json:"statusTavi"`
	Phone             string        `bson:"phone" json:"phone"`
	StartDate         string        `bson:"startDate" json:"startDate"`
}

type ArchivedUser struct {
	ArchivedDate time.Time `bson:"archivedDate" json:"archivedDate"`
	User
}

const (
	admin      = "admin"
	userAdmin  = "user-admin"
	individual = "individual"
	corporate  = "corporate"
)

func (u *User) IsAdmin() bool {
	return u.Role == admin
}

// IsUserManager is true for full admins and user-admins (user lifecycle management).
func (u *User) IsUserManager() bool {
	return u.Role == admin || u.Role == userAdmin
}

func (u *User) GetFullname() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetName() string {
	if u.Role == corporate {
		return u.CorporateName
	}
	return u.GetFullname()
}

func (u *User) GetBankAccountName() string {
	return u.BankAccountName
}

func (u *User) GetThaiCitizenID() string {
	return u.ThaiCitizenID
}

func (u *User) GetAddress() string {
	return u.Address
}
func (u *User) GetRole() string {
	return u.Role
}

func (u *User) GetStatusTavi() bool {
	return u.StatusTavi
}

func (u *User) IsFullnameEmpty() bool {
	return u.FirstName == "" || u.LastName == ""
}

func (u *User) ValidateRole() error {
	if u.Role != corporate && u.Role != individual && u.Role != admin && u.Role != userAdmin {
		return ErrInvalidUserRole
	}
	return nil
}

func (u *User) ValidateVat() error {
	if u.Vat != "N" && u.Vat != "Y" {
		return ErrInvalidUserVat
	}
	return nil
}

func (u *User) ValidateEmail() error {
	if !emailRegexp.MatchString(u.Email) {
		return ErrInvalidFormat
	}
	return nil
}

// ApplyProfileUpdate copies editable profile fields from src onto u.
// Empty strings leave the current value except Role, Vat, Phone, and StatusTavi
// which are always replaced (matching the previous update behavior).
func (u *User) ApplyProfileUpdate(src User) {
	if src.FirstName != "" {
		u.FirstName = toFirstUpper(src.FirstName)
	}
	if src.LastName != "" {
		u.LastName = toFirstUpper(src.LastName)
	}
	if src.CorporateName != "" {
		u.CorporateName = src.CorporateName
	}
	if src.BankAccountName != "" {
		u.BankAccountName = src.BankAccountName
	}
	if src.BankAccountNumber != "" {
		u.BankAccountNumber = extractDigits(src.BankAccountNumber)
	}
	if src.ThaiCitizenID != "" {
		u.ThaiCitizenID = src.ThaiCitizenID
	}
	if src.SiteID != "" {
		u.SiteID = src.SiteID
	}
	if src.Project != "" {
		u.Project = src.Project
	}
	if src.DailyIncome != "" {
		u.DailyIncome = src.DailyIncome
	}
	if src.Address != "" {
		u.Address = src.Address
	}
	if src.StartDate != "" {
		u.StartDate = src.StartDate
	}

	u.StatusTavi = src.StatusTavi
	u.Role = src.Role
	u.Vat = src.Vat
	u.Phone = src.Phone
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func toFirstUpper(s string) string {
	return strings.Title(strings.ToLower(s))
}

func extractDigits(input string) string {
	var result []rune
	for _, char := range input {
		if unicode.IsDigit(char) {
			result = append(result, char)
		}
	}
	return string(result)
}
