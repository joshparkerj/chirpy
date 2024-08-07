package main

import (
	"time"
)

type Token struct {
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
}
