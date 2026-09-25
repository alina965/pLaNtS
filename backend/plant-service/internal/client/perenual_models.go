package client

import (
	"bytes"
	"encoding/json"
)

type SpeciesListResponse struct {
	Data        []Species `json:"data"`
	PerPage     int       `json:"per_page"`
	CurrentPage int       `json:"current_page"`
	LastPage    int       `json:"last_page"`
	Total       int       `json:"total"`
}

type Species struct {
	ID             int      `json:"id"`
	CommonName     string   `json:"common_name"`
	ScientificName []string `json:"scientific_name"`
	OtherName      []string `json:"other_name"`
	DefaultImage   *Image   `json:"default_image"`
}

type SpeciesDetailResponse struct {
	ID                       int                       `json:"id"`
	CommonName               string                    `json:"common_name"`
	ScientificName           []string                  `json:"scientific_name"`
	OtherName                []string                  `json:"other_name"`
	DefaultImage             *Image                    `json:"default_image"`
	GrowthRate               string                    `json:"growth_rate"`
	Type                     string                    `json:"type"`
	Dimensions               FlexibleDimensions        `json:"dimensions"`
	Cycle                    string                    `json:"cycle"`
	Watering                 string                    `json:"watering"`
	WateringGeneralBenchmark *WateringGeneralBenchmark `json:"watering_general_benchmark"`
	Sunlight                 []string                  `json:"sunlight"`
	PruningMonth             []string                  `json:"pruning_month"`
	PruningCount             *PruningCount             `json:"pruning_count"`
	Soil                     []string                  `json:"soil"`
	Maintenance              string                    `json:"maintenance"`
	PoisonousToHumans        bool                      `json:"poisonous_to_humans"`
	PoisonousToPets          bool                      `json:"poisonous_to_pets"`
	CareLevel                string                    `json:"care_level"`
	Description              string                    `json:"description"`
	Flowers                  bool                      `json:"flowers"`
	FloweringSeason          string                    `json:"flowering_season"`
}

type Dimensions struct {
	MinValue float64 `json:"min_value"`
	MaxValue float64 `json:"max_value"`
	Unit     string  `json:"unit"`
}

// FlexibleDimensions accepts object, null, or array from Perenual.
type FlexibleDimensions struct {
	Value *Dimensions
}

func (f *FlexibleDimensions) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" || data[0] == '[' {
		f.Value = nil
		return nil
	}

	var d Dimensions
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	f.Value = &d
	return nil
}

type WateringGeneralBenchmark struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type PruningCount struct {
	Amount   int    `json:"amount"`
	Interval string `json:"interval"`
}

type Image struct {
	RegularUrl string `json:"regular_url"`
	SmallUrl   string `json:"small_url"`
}
