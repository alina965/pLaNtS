package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserPlant struct {
	ID                   uuid.UUID  `json:"id"`
	SpeciesID            int        `json:"speciesId"`
	UserID               uuid.UUID  `json:"userId"`
	Name                 string     `json:"name"`
	WateringIntervalDays int        `json:"wateringIntervalDays"`
	Status               string     `json:"status"`
	LastWateredAt        time.Time  `json:"lastWateredAt"`
	NextWateringAt       time.Time  `json:"nextWateringAt"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            *time.Time `json:"updatedAt"`
}
