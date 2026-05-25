package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Balance      float64   `json:"balance"`
	CreatedAt    time.Time `json:"createdAt"`
}
