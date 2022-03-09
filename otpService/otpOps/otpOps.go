package otpOps

import (
	"context"
	"encoding/json"
	"flamingo-authService/config"
	"flamingo-authService/otpService/twilio"
	"fmt"
)

const (
	otpMessage = "Your OTP is %s . Valid for 15 mins"
)

type OtpOps struct {
	cfg *config.Config
}

func NewOtpOps(cfg *config.Config) *OtpOps {
	return &OtpOps{cfg: cfg}
}

func (o *OtpOps) SendOtp(ctx context.Context, deliveryBody []byte) error {
	var message map[string]interface{}

	if err := json.Unmarshal(deliveryBody, &message); err != nil {
		return err
	}

	otpmessage := fmt.Sprintf(otpMessage, message["otp"].(string))

	twilioModel := twilio.Twilio{
		AccountSid: o.cfg.Twilio.AccountSID,
		AuthToken:  o.cfg.Twilio.AuthToken,
		To:         message["mobile"].(string),
		From:       o.cfg.Twilio.Number,
		Message:    otpmessage,
	}

	status, err := twilioModel.SendOTP()
	if err != nil || status != "success" {
		return err
	}

	return nil
}
