package models

type Token struct {
	Token      string `json:"token"`
	FirstLogin string `json:"firstLogin"`
	User       *User  `json:"user"`
}
