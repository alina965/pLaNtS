package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alina965/pLaNtS/telegram-service/internal/domain"
	"github.com/alina965/pLaNtS/telegram-service/internal/telegram"
	"github.com/google/uuid"
)

type TelegramHandler struct {
	service *telegram.Service
}

func NewTelegramHandler(service *telegram.Service) *TelegramHandler {
	return &TelegramHandler{service: service}
}

func (h *TelegramHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := h.service.CreateLink(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TelegramHandler) Notify(w http.ResponseWriter, r *http.Request) {
	var req domain.NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	if err := h.service.Notify(r.Context(), userID, req.Text); err != nil {
		if errors.Is(err, telegram.ErrNotLinked) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
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
