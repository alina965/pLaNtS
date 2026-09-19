package repository

import (
	"context"
	"time"

	"github.com/alina965/pLaNtS/user-plants-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WateringEventsRepository struct {
	db *pgxpool.Pool
}

func NewWateringEventsRepository(db *pgxpool.Pool) *WateringEventsRepository {
	return &WateringEventsRepository{db: db}
}

func (r *WateringEventsRepository) CreateWateringEvent(event *domain.WateringEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, `
	INSERT INTO watering_events (
		id,
		user_plant_id,
		watered_at,
		created_at
	) VALUES ($1, $2, $3, $4)`,
		event.ID,
		event.UserPlantID,
		event.WateredAt,
		event.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *WateringEventsRepository) GetWateringEventsByUserPlantID(userPlantID uuid.UUID) ([]*domain.WateringEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT id, user_plant_id, watered_at, created_at
		FROM watering_events
		WHERE user_plant_id = $1
		ORDER BY watered_at DESC`, userPlantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*domain.WateringEvent, 0)
	for rows.Next() {
		event := &domain.WateringEvent{}
		if err := rows.Scan(
			&event.ID,
			&event.UserPlantID,
			&event.WateredAt,
			&event.CreatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}
