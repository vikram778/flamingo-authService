package auth

import (
	"context"
	"flamingo-authService/auth/models"
)

type DBOps interface {
	CreateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error)
	GetProfile(ctx context.Context, mobile string) (*models.Profile, error)
	UpdateProfileStatus(ctx context.Context, profile *models.Profile) error
	CreateOTP(ctx context.Context, otp *models.OTP) (*models.OTP, error)
	VerifyOTP(ctx context.Context, otp string, profileID int64) (bool, int64, error)
	UpdateOtpStatus(ctx context.Context, status bool, id int64) error
}
