package grpc

import (
	"context"
	"encoding/json"
	"flamingo-authService/auth/models"
	authService "flamingo-authService/auth/proto"
	"flamingo-authService/auth/rabbitmq"
	"flamingo-authService/auth/repository"
	"flamingo-authService/auth/utils"
	"flamingo-authService/config"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

// AuthMicroservice gRPC microservice
type AuthMicroservice struct {
	cfg     *config.Config
	repo    repository.DBOps
	publish rabbitmq.Publisher
}

// NewAuthMicroservice gRPC microservice constructor
func NewAuthMicroservice(repo repository.DBOps, cfg *config.Config, publisher rabbitmq.Publisher) *AuthMicroservice {
	return &AuthMicroservice{repo: repo, cfg: cfg, publish: publisher}
}

//Signup handles signup request creates profile and generates OTP and saves OTP in otp_log table and publishes OTP to rabbitmq
func (a *AuthMicroservice) Signup(ctx context.Context, req *authService.SignUpRequest) (*authService.SignUpResponse, error) {

	var (
		err     error
		message map[string]interface{}
	)

	profile := &models.Profile{
		Name:     req.Name,
		Mobile:   req.Mobile,
		DOB:      req.Dob,
		Location: req.Location,
	}

	if profile, err = a.repo.CreateProfile(ctx, profile); err != nil {
		return nil, err
	}

	message["otp"] = utils.GenerateOTPCode(6)
	message["mobile"] = req.Mobile

	otp := &models.OTP{
		ProfileID: profile.ID,
		OTP:       message["otp"].(string),
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(15 * time.Minute),
	}

	if _, err = a.repo.CreateOTP(ctx, otp); err != nil {
		return nil, err
	}

	msg, _ := json.Marshal(message)
	if err := a.publish.Publish(msg); err != nil {
		return nil, err
	}
	return &authService.SignUpResponse{Status: "Ok"}, nil
}

//VerifyOtp verifies OTP and validates login and signup
func (a *AuthMicroservice) VerifyOtp(ctx context.Context, req *authService.VerifyOtpRequest) (*authService.VerifyOtpResponse, error) {

	profile, err := a.repo.GetProfile(ctx, req.Mobile)
	if err != nil {
		return nil, err
	}

	status, id, err := a.repo.VerifyOTP(ctx, strconv.Itoa(int(req.Otp)), profile.ID)
	if err != nil {
		return nil, err
	}

	err = a.repo.UpdateOtpStatus(ctx, true, id)
	if err != nil {
		return nil, err
	}

	if !status {
		return &authService.VerifyOtpResponse{Status: "fail"}, nil
	}

	profile.IsVerified = true
	if err = a.repo.UpdateProfileStatus(ctx, profile); err != nil {
		return nil, err
	}

	return &authService.VerifyOtpResponse{Status: "success"}, nil
}

//Login generates otp and creates otp entry in otp_log table and publishes generated otp to rabbitmq to be consumed by otp service
func (a *AuthMicroservice) Login(ctx context.Context, req *authService.LoginRequest) (*authService.LoginResponse, error) {
	var (
		err     error
		message map[string]interface{}
	)

	profile, err := a.repo.GetProfile(ctx, req.Mobile)
	if err != nil {
		return nil, err
	}

	message["otp"] = utils.GenerateOTPCode(6)
	message["mobile"] = req.Mobile

	otp := &models.OTP{
		ProfileID: profile.ID,
		OTP:       message["otp"].(string),
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(15 * time.Minute),
	}

	if _, err = a.repo.CreateOTP(ctx, otp); err != nil {
		return nil, err
	}

	msg, _ := json.Marshal(message)
	if err := a.publish.Publish(msg); err != nil {
		return nil, err
	}
	return &authService.LoginResponse{Status: "otp sent"}, nil
}

//GetProfile fetches profile of the user
func (a *AuthMicroservice) GetProfile(ctx context.Context, req *authService.GetProfileRequest) (*authService.GetProfileResponse, error) {
	profile, err := a.repo.GetProfile(ctx, req.Mobile)
	if err != nil {
		return nil, err
	}
	return &authService.GetProfileResponse{Profile: a.convertEmailToProto(profile)}, nil
}

func (a *AuthMicroservice) convertEmailToProto(profile *models.Profile) *authService.Profile {
	return &authService.Profile{
		Name:      profile.Name,
		Mobile:    profile.Mobile,
		Dob:       profile.DOB,
		Location:  profile.Location,
		CreatedAt: timestamppb.New(profile.CreatedAt),
	}
}
