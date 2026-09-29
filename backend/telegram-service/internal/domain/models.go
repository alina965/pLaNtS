package domain

import (
	"time"

	"github.com/google/uuid"
)

type TelegramLink struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	ChatID    *int64     `json:"chatId"`
	LinkCode  *string    `json:"linkCode"`
	CreatedAt time.Time  `json:"createdAt"`
	LinkedAt  *time.Time `json:"linkedAt"`
}
