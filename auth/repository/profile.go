package repository

import (
	"context"
	"flamingo-authService/auth/models"
	"time"
)

func (r *DBOps) CreateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error) {

	profile.CreatedAt = time.Now()

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		createProfileQuery,
		profile.Name,
		profile.Mobile,
		profile.DOB,
		profile.Location,
		profile.IsVerified,
		profile.CreatedAt,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	profile.ID = id
	return profile, nil
}

func (r *DBOps) GetProfile(ctx context.Context, mobile string) (*models.Profile, error) {
	var profile models.Profile

	if err := r.db.Get(&profile, getProfileQuery, mobile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *DBOps) UpdateProfileStatus(ctx context.Context, profile *models.Profile) error {

	if _, err := r.db.Exec(updateProfileStatus, profile.IsVerified, profile.Mobile); err != nil {
		return err
	}
	return nil
}
