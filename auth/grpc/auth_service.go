package grpc

import (
	"context"
	"flamingo-authService/auth/models"
	authService "flamingo-authService/auth/proto"
	"flamingo-authService/auth/rabbitmq"
	"flamingo-authService/auth/repository"
	"flamingo-authService/config"
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

func (a *AuthMicroservice) Signup(ctx context.Context, req *authService.SignUpRequest) (*authService.SignUpResponse, error) {

	profile := &models.Profile{
		Name:     req.Name,
		Mobile:   req.Mobile,
		DOB:      req.Dob,
		Location: req.Location,
	}

	if _, err := a.repo.CreateProfile(ctx, profile); err != nil {
		return nil, err
	}

	msg := "sendOTP"

	if err := a.publish.Publish([]byte(msg)); err != nil {
		return nil, err
	}
	return &authService.SignUpResponse{Status: "Ok"}, nil
}

func (a *AuthMicroservice) VerifyOtp(ctx context.Context, req *authService.VerifyOtpRequest) (*authService.VerifyOtpResponse, error) {

	panic("not implemented")
}

func (a *AuthMicroservice) Login(ctx context.Context, req *authService.LoginRequest) (*authService.LoginResponse, error) {
	panic("not implemented")
}

func (a *AuthMicroservice) GetProfile(ctx context.Context, req *authService.GetProfileRequest) (*authService.GetProfileResponse, error) {
	panic("not implemented")
}
