package service

import (
	"context"
	"fmt"
	"time"

	"github.com/fdg312/ecommerce-microservices/internal/models"
	"github.com/fdg312/ecommerce-microservices/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(r *repository.UserRepository) *AuthService {
	return &AuthService{r}
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("couldn`t hash password: %w", err)
	}
	passwordHash := string(hashBytes)

	var user models.User

	user.ID = uuid.New()
	user.Email = email
	user.PasswordHash = passwordHash

	err = s.repo.CreateUser(ctx, user)

	if err != nil {
		return fmt.Errorf("coudln`t create user: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("couldn`t get user by email: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials: %w", err)
	}

	var secretKey = []byte("my-super-secret-key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	readyToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("couldn`t sign token: %w", err)
	}

	return readyToken, nil
}
