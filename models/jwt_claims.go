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

func (u *UserClaims) GetStatusTavi() bool {
	return u.StatusTavi
}
