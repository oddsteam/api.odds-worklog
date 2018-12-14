package models

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	User *User `json:"user"`
	jwt.StandardClaims
}
