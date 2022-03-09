package otpService

import "context"

type OtpOps interface {
	SendOtp(ctx context.Context, deliveryBody []byte) error
}
