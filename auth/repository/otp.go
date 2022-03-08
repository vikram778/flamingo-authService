package repository

import (
	"context"
	"flamingo-authService/auth/models"
	"time"
)

func (r *DBOps) CreateOTP(ctx context.Context, otp *models.OTP) (*models.OTP, error) {
	return &models.OTP{}, nil
}

func (r *DBOps) VerifyOTP(ctx context.Context, otp, profileID int64) (bool, int64, error) {
	var otpModel models.OTP

	if err := r.db.Get(&otpModel, verifyOTPQuery, profileID, time.Now()); err != nil || otpModel.Validated == true {
		return false, 0, err
	}

	if otp == otpModel.OTP {
		return true, otpModel.ID, nil
	}

	return false, 0, nil
}

func (r *DBOps) UpdateOtpStatus(ctx context.Context, status bool, id int64) error {

	if _, err := r.db.Exec(updateOTPStatusQuery, status, id); err != nil {
		return err
	}
	return nil
}
