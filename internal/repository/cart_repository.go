package repository

import (
	"context"
	"encoding/json"

	"github.com/fdg312/ecommerce-microservices/pkg/cart_v1"
	"github.com/redis/go-redis/v9"
)

type CartRepository struct {
	client *redis.Client
}

func NewCartRepository(client *redis.Client) *CartRepository {
	return &CartRepository{client}
}

func (r *CartRepository) AddItem(ctx context.Context, userID string, item *cart_v1.AddItemRequest) error {
	key := "cart-" + userID

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, key, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
