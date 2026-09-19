package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserPlant struct {
	ID                   uuid.UUID  `db:"id"`
	SpeciesID            int        `db:"species_id"`
	UserID               uuid.UUID  `db:"user_id"`
	Name                 string     `db:"name"`
	WateringIntervalDays int        `db:"watering_interval_days"`
	Status               string     `db:"status"`
	LastWateredAt        time.Time  `db:"last_watered_at"`
	NextWateringAt       time.Time  `db:"next_watering_at"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            *time.Time `db:"updated_at"`
}

type WateringEvent struct {
	ID          uuid.UUID `db:"id"`
	UserPlantID uuid.UUID `db:"user_plant_id"`
	WateredAt   time.Time `db:"watered_at"`
	CreatedAt   time.Time `db:"created_at"`
}
