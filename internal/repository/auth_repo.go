package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"tt/internal/models"
)

type AuthRepositoryInterface interface {
	CreateUser(ctx context.Context, user models.User) (uint64, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepositoryInterface {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) (uint64, error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`

	var id int64
	err := r.db.QueryRow(ctx, query, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
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
