package services

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"tt/internal/models"
	"tt/internal/repository"
	"tt/pkg/utils"
)

type AuthServiceInterface interface {
}

type AuthService struct {
	repo repository.AuthRepositoryInterface
	rdb  *redis.Client
}

func (s *AuthService) Register(ctx context.Context, user models.User) (uint64, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	user.Username = hashedPassword

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (s *AuthService) Login(ctx context.Context, user models.User) (*string, error) {
	userDB, err := s.repo.GetUser(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %w", err)
	}

	if !utils.VerifyPassword(user.Password, userDB.Password) {
		return nil, fmt.Errorf("invalid password for user: %s", user.Username)
	}

	token, err := utils.GenerateJWT(userDB.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	err = s.saveTokenToRedis(ctx, token, strconv.FormatUint(userDB.ID, 10))
	if err != nil {
		return nil, fmt.Errorf("failed to save JWT token in redis: %w", err)
	}
	return &token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string, userID string) error {
	key := fmt.Sprintf("auth:%s:%s", token, userID)

	cmd := s.rdb.Del(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (s *AuthService) saveTokenToRedis(ctx context.Context, token string, userID string) error {
	ttl := 30 * time.Minute

	key := fmt.Sprintf("auth:%s:%s", token, userID)

	err := s.rdb.Set(ctx, key, token, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save token to Redis: %w", err)
	}

	return nil
}

func (s *AuthService) GetTokenFromRedis(ctx context.Context, token string, userID string) (string, error) {
	key := fmt.Sprintf("auth:%s:%s", token, userID)
	token, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("token not found or expired")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get token from Redis: %w", err)
	}

	return token, nil
}
