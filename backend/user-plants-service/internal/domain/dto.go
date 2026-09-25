package domain

import "time"

type UserPlantsResponse struct {
	UserPlantsList []*UserPlant `json:"userPlantsList"`
}

type CreateUserPlantRequest struct {
	SpeciesID            int       `json:"speciesId"`
	Name                 string    `json:"name"`
	WateringIntervalDays int       `json:"wateringIntervalDays"`
	LastWateredAt        time.Time `json:"lastWateredAt"`
}

type UpdateUserPlantRequest struct {
	Name                 *string    `json:"name"`
	WateringIntervalDays *int       `json:"wateringIntervalDays"`
	Status               *string    `json:"status"`
	NextWateringAt       *time.Time `json:"nextWateringAt"`
}

type MarkWateredRequest struct {
	WateredAt *time.Time `json:"wateredAt"`
}

type UpdateStatusRequest struct {
	Status *string `json:"status"`
}
