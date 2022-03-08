package models

import (
	"time"
)

type Profile struct {
	ID         int64     `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Mobile     string    `json:"mobile" db:"mobile"`
	DOB        string    `json:"dob" db:"dob"`
	Location   string    `json:"location" db:"location"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
