package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alina965/pLaNtS/user-plants-service/internal/domain"
	"github.com/alina965/pLaNtS/user-plants-service/internal/repository"
	"github.com/alina965/pLaNtS/user-plants-service/internal/user_plant"
	"github.com/google/uuid"
)

type UserPlantsHandler struct {
	service *user_plant.UserPlantsService
}

func NewUserPlantsHandler(service *user_plant.UserPlantsService) *UserPlantsHandler {
	return &UserPlantsHandler{service: service}
}

func (h *UserPlantsHandler) CreateUserPlant(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req domain.CreateUserPlantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plant, err := h.service.CreatePlant(r.Context(), req.SpeciesID, userID, req.Name, req.WateringIntervalDays, req.LastWateredAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, plant)
}

func (h *UserPlantsHandler) GetUserPlants(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	plants, err := h.service.GetPlantsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, plants)
}

func (h *UserPlantsHandler) GetUserPlantByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid plant id", http.StatusBadRequest)
		return
	}

	plant, err := h.service.GetPlantByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if plant == nil || plant.UserID != userID {
		http.Error(w, repository.ErrUserPlantNotFound.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, plant)
}

func (h *UserPlantsHandler) UpdateUserPlant(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid plant id", http.StatusBadRequest)
		return
	}

	var req domain.UpdateUserPlantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdatePlant(r.Context(), id, userID, req.Name, req.WateringIntervalDays, req.Status, req.NextWateringAt)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserPlantsHandler) DeleteUserPlant(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid plant id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePlant(r.Context(), id, userID); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserPlantsHandler) MarkWatered(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid plant id", http.StatusBadRequest)
		return
	}

	var req domain.MarkWateredRequest
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	plant, err := h.service.GetPlantByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if plant == nil || plant.UserID != userID {
		http.Error(w, repository.ErrUserPlantNotFound.Error(), http.StatusNotFound)
		return
	}

	wateredAt := time.Now()
	if req.WateredAt != nil {
		wateredAt = *req.WateredAt
	}
	next := wateredAt.AddDate(0, 0, plant.WateringIntervalDays)

	if err := h.service.MarkWatered(r.Context(), id, userID, wateredAt, next); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserPlantsHandler) GetWateringEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid plant id", http.StatusBadRequest)
		return
	}

	events, err := h.service.GetWateringEvents(r.Context(), id, userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, events)
}

func (h *UserPlantsHandler) ListNeedingWater(w http.ResponseWriter, r *http.Request) {
	plants, err := h.service.ListNeedingWater(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, plants)
}

func (h *UserPlantsHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid plant id", http.StatusBadRequest)
		return
	}

	var req domain.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Status == nil {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}
	err = h.service.UpdateStatus(r.Context(), id, *req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func userIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, false
	}

	return userID, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrUserPlantNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Error(w, err.Error(), http.StatusBadRequest)
}
