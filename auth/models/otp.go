package models

import "time"

type OTP struct {
	ID        int64     `json:"id" db:"id"`
	ProfileID int64     `json:"profile_id" db:"profile_id"`
	OTP       string    `json:"otp" db:"otp"`
	Validated bool      `json:"validated" db:"validated"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Expiry    time.Time `json:"expiry" db:"expiry"`
}
