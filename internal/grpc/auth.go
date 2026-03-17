package grpc

import (
	"context"
	"fmt"

	"github.com/fdg312/ecommerce-microservices/internal/service"
	"github.com/fdg312/ecommerce-microservices/pkg/auth_v1"
)

type AuthServer struct {
	service *service.AuthService
	auth_v1.UnimplementedAuthServiceServer
}

func NewAuthServer(s *service.AuthService) *AuthServer {
	return &AuthServer{service: s}
}

func (s *AuthServer) Register(ctx context.Context, r *auth_v1.RegisterRequest) (*auth_v1.RegisterResponse, error) {
	err := s.service.Register(ctx, r.GetEmail(), r.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("couldn`t : %w", err)
	}
	return &auth_v1.RegisterResponse{UserId: "some-id"}, nil
}

func (s *AuthServer) Login(ctx context.Context, r *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	token, err := s.service.Login(ctx, r.GetEmail(), r.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("couldn`t : %w", err)
	}
	return &auth_v1.LoginResponse{Token: token}, nil
}

func (s *AuthServer) VerifyToken(ctx context.Context, r *auth_v1.VerifyTokenRequest) (*auth_v1.VerifyTokenResponse, error) {
	userId, err := s.service.VerifyToken(ctx, r.GetToken())
	if err != nil {
		return nil, fmt.Errorf("couldn`t : %w", err)
	}
	return &auth_v1.VerifyTokenResponse{UserId: userId}, nil
}
