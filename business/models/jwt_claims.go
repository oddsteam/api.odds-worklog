package models

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	User *UserClaims `json:"user"`
	jwt.StandardClaims
}

type UserClaims struct {
	ID         string `json:"id,omitempty"`
	Role       string `json:"role"`
	StatusTavi bool   `json:"statusTavi"`
}

func (u *UserClaims) IsAdmin() bool {
	return u.Role == admin
}

// IsUserManager is true for JWTs issued to full admins or user-admins.
func (u *UserClaims) IsUserManager() bool {
	return u.Role == admin || u.Role == userAdmin
}

// CanExportIncome is true only for full admins (income export endpoints).
func (u *UserClaims) CanExportIncome() bool {
	return u.IsAdmin()
}

func (u *UserClaims) GetStatusTavi() bool {
	return u.StatusTavi
}
