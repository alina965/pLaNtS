package repository

import (
	"context"
	"errors"
	"time"

	"github.com/alina965/pLaNtS/telegram-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTelegramLinkNotFound = errors.New("telegram link not found")

const telegramLinkColumns = `
	id, user_id, chat_id, link_code, created_at, linked_at`

type TelegramLinksRepository struct {
	db *pgxpool.Pool
}

func NewTelegramLinksRepository(db *pgxpool.Pool) *TelegramLinksRepository {
	return &TelegramLinksRepository{db: db}
}

func scanTelegramLink(row pgx.Row) (*domain.TelegramLink, error) {
	link := &domain.TelegramLink{}
	err := row.Scan(
		&link.ID,
		&link.UserID,
		&link.ChatID,
		&link.LinkCode,
		&link.CreatedAt,
		&link.LinkedAt,
	)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *TelegramLinksRepository) CreateLink(ctx context.Context, link *domain.TelegramLink) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO telegram_links (
			id,
			user_id,
			chat_id,
			link_code,
			created_at,
			linked_at
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		link.ID,
		link.UserID,
		link.ChatID,
		link.LinkCode,
		link.CreatedAt,
		link.LinkedAt,
	)
	return err
}

func (r *TelegramLinksRepository) UpsertLinkCode(ctx context.Context, userID uuid.UUID, linkCode string) (*domain.TelegramLink, error) {
	now := time.Now().UTC()
	row := r.db.QueryRow(ctx, `
		INSERT INTO telegram_links (id, user_id, chat_id, link_code, created_at, linked_at)
		VALUES ($1, $2, NULL, $3, $4, NULL)
		ON CONFLICT (user_id) DO UPDATE
		SET link_code = EXCLUDED.link_code,
		    created_at = EXCLUDED.created_at
		RETURNING`+telegramLinkColumns,
		uuid.New(),
		userID,
		linkCode,
		now,
	)

	link, err := scanTelegramLink(row)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *TelegramLinksRepository) FindByLinkCode(ctx context.Context, linkCode string) (*domain.TelegramLink, error) {
	row := r.db.QueryRow(ctx, `
		SELECT`+telegramLinkColumns+`
		FROM telegram_links
		WHERE link_code = $1`, linkCode)

	link, err := scanTelegramLink(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *TelegramLinksRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.TelegramLink, error) {
	row := r.db.QueryRow(ctx, `
		SELECT`+telegramLinkColumns+`
		FROM telegram_links
		WHERE user_id = $1`, userID)

	link, err := scanTelegramLink(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *TelegramLinksRepository) AttachChat(ctx context.Context, linkCode string, chatID int64) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(ctx, `
		UPDATE telegram_links
		SET
			chat_id = $1,
			linked_at = $2,
			link_code = NULL
		WHERE link_code = $3`,
		chatID,
		now,
		linkCode,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrTelegramLinkNotFound
	}
	return nil
}
