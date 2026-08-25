package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alina965/pLaNtS/plant-service/internal/domain"
)

type PlantsService interface {
	GetSpecies(page int, query string) (*domain.SpeciesList, error)
	GetSpeciesDetails(speciesID int) (*domain.SpeciesDetails, error)
}

type PlantsHandler struct {
	service PlantsService
}

func NewPlantsHandler(service PlantsService) *PlantsHandler {
	return &PlantsHandler{service: service}
}

func (h *PlantsHandler) GetSpeciesList(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")

	var page int
	if pageStr == "" {
		page = 1
	} else {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if p < 1 {
			p = 1
		}

		page = p
	}

	query := r.URL.Query().Get("query")

	species, err := h.service.GetSpecies(page, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(species)
}

func (h *PlantsHandler) GetSpeciesDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	speciesDetails, err := h.service.GetSpeciesDetails(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(speciesDetails)
}
