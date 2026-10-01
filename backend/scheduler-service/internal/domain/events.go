package domain

import (
	"time"

	"github.com/google/uuid"
)

type WateringNotifyEvent struct {
	UserID     uuid.UUID `json:"user_id"`
	PlantID    uuid.UUID `json:"plant_id"`
	PlantName  string    `json:"plant_name"`
	Status     string    `json:"status"`
	Text       string    `json:"text"`
	OccurredAt time.Time `json:"occurred_at"`
}
