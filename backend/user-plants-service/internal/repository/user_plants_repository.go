package repository

import (
	"context"
	"errors"
	"time"

	"github.com/alina965/pLaNtS/user-plants-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserPlantNotFound = errors.New("user plant not found")

const userPlantColumns = `
	id, species_id, user_id, name, watering_interval_days,
	status, last_watered_at, next_watering_at, created_at, updated_at`

type UserPlantsRepository struct {
	db *pgxpool.Pool
}

func NewUserPlantsRepository(db *pgxpool.Pool) *UserPlantsRepository {
	return &UserPlantsRepository{db: db}
}

func scanUserPlant(row pgx.Row) (*domain.UserPlant, error) {
	plant := &domain.UserPlant{}
	err := row.Scan(
		&plant.ID,
		&plant.SpeciesID,
		&plant.UserID,
		&plant.Name,
		&plant.WateringIntervalDays,
		&plant.Status,
		&plant.LastWateredAt,
		&plant.NextWateringAt,
		&plant.CreatedAt,
		&plant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return plant, nil
}

func (r *UserPlantsRepository) CreateUserPlant(plant *domain.UserPlant) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, `
	INSERT INTO user_plants (
		id,
		species_id,
		user_id,
		name,
		watering_interval_days,
		status,
		last_watered_at,
		next_watering_at,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		plant.ID,
		plant.SpeciesID,
		plant.UserID,
		plant.Name,
		plant.WateringIntervalDays,
		plant.Status,
		plant.LastWateredAt,
		plant.NextWateringAt,
		plant.CreatedAt,
		plant.UpdatedAt,
	)
	return err
}

func (r *UserPlantsRepository) GetUserPlantByID(id uuid.UUID) (*domain.UserPlant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRow(ctx, `
		SELECT`+userPlantColumns+`
		FROM user_plants
		WHERE id = $1`, id)

	plant, err := scanUserPlant(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return plant, nil
}

func (r *UserPlantsRepository) GetUserPlantsByUserID(userID uuid.UUID) ([]*domain.UserPlant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT`+userPlantColumns+`
		FROM user_plants
		WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userPlants := make([]*domain.UserPlant, 0)
	for rows.Next() {
		plant, err := scanUserPlant(rows)
		if err != nil {
			return nil, err
		}
		userPlants = append(userPlants, plant)
	}

	return userPlants, rows.Err()
}

func (r *UserPlantsRepository) UpdateUserPlant(plant *domain.UserPlant) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.db.Exec(ctx, `
		UPDATE user_plants
		SET
			name = $1,
			watering_interval_days = $2,
			status = $3,
			next_watering_at = $4,
			updated_at = $5
		WHERE id = $6 AND user_id = $7`,
		plant.Name,
		plant.WateringIntervalDays,
		plant.Status,
		plant.NextWateringAt,
		plant.UpdatedAt,
		plant.ID,
		plant.UserID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrUserPlantNotFound
	}

	return nil
}

func (r *UserPlantsRepository) DeleteUserPlantByID(id uuid.UUID, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.db.Exec(ctx, `
		DELETE FROM user_plants
		WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrUserPlantNotFound
	}

	return nil
}

func (r *UserPlantsRepository) MarkWatered(plant *domain.UserPlant) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.db.Exec(ctx, `
		UPDATE user_plants
		SET
			status = $1,
			next_watering_at = $2,
			last_watered_at = $3,
			updated_at = $4
		WHERE id = $5 AND user_id = $6`,
		plant.Status,
		plant.NextWateringAt,
		plant.LastWateredAt,
		plant.UpdatedAt,
		plant.ID,
		plant.UserID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrUserPlantNotFound
	}

	return nil
}
