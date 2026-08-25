package repository

import (
	"context"
	"errors"
	"time"

	"github.com/alina965/pLaNtS/user-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(token *domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, "INSERT INTO refresh_tokens (id, token, expires_at, is_revoked, user_id, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		token.ID, token.Token, token.ExpiresAt, token.IsRevoked, token.UserID, time.Now())

	return err
}

func (r *RefreshTokenRepository) FindRefreshTokenByToken(token string) (*domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRow(ctx, "SELECT * FROM refresh_tokens WHERE token = $1", token)

	var refreshToken domain.RefreshToken
	err := row.Scan(&refreshToken.ID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.IsRevoked, &refreshToken.UserID, &refreshToken.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) DeleteRefreshTokenById(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, "DELETE FROM refresh_tokens WHERE id = $1", id)

	return err
}
