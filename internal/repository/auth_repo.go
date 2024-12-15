package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"tt/internal/models"
)

type AuthRepositoryInterface interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, username string) (models.User, error)
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, query, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetUser(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, password, created_at FROM users WHERE username = $1`

	row := r.db.QueryRow(ctx, query, username)

	user := models.User{}

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, err
}
