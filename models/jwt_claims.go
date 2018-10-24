package models

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}
