package repository

import (
	"context"
	"fmt"

	"github.com/fdg312/ecommerce-microservices/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	query := "INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("сouldn't create a user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	query := "SELECT id, email, password_hash FROM users WHERE email = $1"
	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return user, fmt.Errorf("couldn't get user: %w", err)
	}
	return user, nil
}
