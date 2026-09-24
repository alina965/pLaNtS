package user_plant

import (
	"context"
	"errors"
	"time"

	"github.com/alina965/pLaNtS/user-plants-service/internal/domain"
	"github.com/alina965/pLaNtS/user-plants-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPlantsService struct {
	plantRepo  *repository.UserPlantsRepository
	eventsRepo *repository.WateringEventsRepository
	db         *pgxpool.Pool
}

func NewUserPlantsService(plantRepo *repository.UserPlantsRepository, eventsRepo *repository.WateringEventsRepository, db *pgxpool.Pool) *UserPlantsService {
	return &UserPlantsService{plantRepo: plantRepo, eventsRepo: eventsRepo, db: db}
}

func (s *UserPlantsService) CreatePlant(speciesID int, userID uuid.UUID, name string, interval int, last time.Time) (*domain.UserPlant, error) {
	if speciesID <= 0 {
		return nil, errors.New("invalid speciesID")
	}
	if userID == uuid.Nil {
		return nil, errors.New("invalid userID")
	}
	if interval <= 0 {
		return nil, errors.New("invalid interval")
	}
	if last.IsZero() || last.After(time.Now()) {
		return nil, errors.New("invalid time")
	}

	next := last.AddDate(0, 0, interval)

	now := time.Now()
	plant := domain.UserPlant{
		ID:                   uuid.New(),
		SpeciesID:            speciesID,
		UserID:               userID,
		Name:                 name,
		WateringIntervalDays: interval,
		Status:               "ok",
		LastWateredAt:        last,
		NextWateringAt:       next,
		CreatedAt:            now,
		UpdatedAt:            &now,
	}

	err := s.plantRepo.CreateUserPlant(&plant)
	if err != nil {
		return nil, err
	}

	return &plant, nil
}

func (s *UserPlantsService) GetPlantByID(id uuid.UUID) (*domain.UserPlant, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid ID")
	}

	plant, err := s.plantRepo.GetUserPlantByID(id)
	if err != nil {
		return nil, err
	}

	return plant, nil
}

func (s *UserPlantsService) GetPlantsByUserID(userID uuid.UUID) (*domain.UserPlantsResponse, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	plants, err := s.plantRepo.GetUserPlantsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserPlantsResponse{UserPlantsList: plants}, nil
}

func (s *UserPlantsService) UpdatePlant(id uuid.UUID, userID uuid.UUID, name *string, interval *int, status *string, nextWatering *time.Time) error {
	if userID == uuid.Nil || id == uuid.Nil {
		return errors.New("invalid id")
	}

	if status == nil || (*status != "ok" && *status != "due" && *status != "overdue") {
	}

	plant, err := s.plantRepo.GetUserPlantByID(id)
	if err != nil {
		return err
	}
	if plant == nil || plant.UserID != userID {
		return repository.ErrUserPlantNotFound
	}

	if name != nil {
		plant.Name = *name
	}
	if interval != nil {
		plant.WateringIntervalDays = *interval
	}
	if status != nil {
		plant.Status = *status
	}
	if nextWatering != nil {
		plant.NextWateringAt = *nextWatering
	}

	return s.plantRepo.UpdateUserPlant(plant)
}

func (s *UserPlantsService) DeletePlant(id uuid.UUID, userID uuid.UUID) error {
	if id == uuid.Nil || userID == uuid.Nil {
		return errors.New("invalid ID")
	}

	err := s.plantRepo.DeleteUserPlantByID(id, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserPlantsService) MarkWatered(id uuid.UUID, userID uuid.UUID, last time.Time, next time.Time) error {
	if id == uuid.Nil || userID == uuid.Nil {
		return errors.New("invalid ID")
	}
	if last.IsZero() || next.IsZero() || last.After(next) || last.After(time.Now()) {
		return errors.New("invalid time")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	plant := domain.UserPlant{ID: id, UserID: userID, LastWateredAt: last, NextWateringAt: next, Status: "ok"}
	if err = s.plantRepo.MarkWatered(ctx, tx, &plant); err != nil {
		return err
	}

	event := domain.WateringEvent{
		ID:          uuid.New(),
		UserPlantID: id,
		WateredAt:   last,
		CreatedAt:   time.Now(),
	}
	if err = s.eventsRepo.CreateWateringEvent(ctx, tx, &event); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *UserPlantsService) GetWateringEvents(plantID uuid.UUID, userID uuid.UUID) ([]*domain.WateringEvent, error) {
	plant, err := s.plantRepo.GetUserPlantByID(plantID)
	if err != nil {
		return nil, err
	}

	if plant == nil || plant.UserID != userID {
		return nil, repository.ErrUserPlantNotFound
	}

	return s.eventsRepo.GetWateringEventsByUserPlantID(plantID)
}
