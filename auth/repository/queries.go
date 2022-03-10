package repository

const (
	createProfileQuery  = `INSERT INTO profile (name, mobile, dob, location, is_verified, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	getProfileQuery     = `SELECT * FROM profile WHERE mobile = $1`
	updateProfileStatus = `UPDATE profile SET is_verified = $1 WHERE mobile = $2`

	createOTPQuery       = `INSERT INTO otp_log (profile_id, otp, validated, created_at, expiry) VALUES ($1, $2, $3, $4, $5)`
	verifyOTPQuery       = `SELECT * FROM otp_log WHERE profile_id = $1 AND validated = 'f' AND expiry > $2 LIMIT 1`
	updateOTPStatusQuery = `UPDATE otp_log SET validated = $1 WHERE id = $2`
)
