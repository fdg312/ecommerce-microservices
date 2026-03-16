package service

import (
	"context"

	"github.com/fdg312/ecommerce-microservices/internal/repository"
	"github.com/fdg312/ecommerce-microservices/pkg/auth_v1"
	"github.com/fdg312/ecommerce-microservices/pkg/cart_v1"
)

type CartService struct {
	client auth_v1.AuthServiceClient
	repo   *repository.CartRepository
	cart_v1.UnimplementedCartServiceServer
}

func NewCartService(c auth_v1.AuthServiceClient, r *repository.CartRepository) *CartService {
	return &CartService{client: c, repo: r}
}

func (s *CartService) AddItem(ctx context.Context, req *cart_v1.AddItemRequest) (*cart_v1.AddItemResponse, error) {
	verifyResp, err := s.client.VerifyToken(ctx, &auth_v1.VerifyTokenRequest{})
	if err != nil {
		return nil, err
	}

	userID := verifyResp.GetUserId()

	err = s.repo.AddItem(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return &cart_v1.AddItemResponse{ItemId: "123"}, nil
}
