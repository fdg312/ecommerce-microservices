package grpc

import (
	"context"
	"strings"

	"github.com/fdg312/ecommerce-microservices/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type InterceptorManager struct {
	service service.AuthService
}

func NewInterceptorManager(s service.AuthService) *InterceptorManager {
	return &InterceptorManager{s}
}

func (i *InterceptorManager) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if info.FullMethod == "/auth_v1.AuthService/Login" || info.FullMethod == "/auth_v1.AuthService/Register" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	if len(md.Get("authorization")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization is not provided")
	}

	trimToken := strings.Trim(md.Get("authorization")[0], " ")
	if !strings.HasPrefix(trimToken, "Bearer ") {
		return nil, status.Errorf(codes.Unauthenticated, "authorization is not provided")
	}
	pureToken := strings.TrimPrefix(trimToken, "Bearer ")
	id, err := i.service.VerifyToken(ctx, pureToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization is not provided")
	}

	ctx = context.WithValue(ctx, UserIDKey, id)

	return handler(ctx, req)
}
