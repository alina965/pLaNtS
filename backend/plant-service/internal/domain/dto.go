package domain

type SpeciesList struct {
	Species     []Species `json:"species"`
	PerPage     int       `json:"per_page"`
	CurrentPage int       `json:"current_page"`
	LastPage    int       `json:"last_page"`
	Total       int       `json:"total"`
}

type SpeciesDetails struct {
	ID                       int                       `json:"id"`
	CommonName               string                    `json:"common_name"`
	ScientificNames          []string                  `json:"scientific_names"`
	OtherNames               []string                  `json:"other_names"`
	Image                    *Image                    `json:"image"`
	GrowthRate               string                    `json:"growth_rate"`
	Type                     string                    `json:"type"`
	Dimensions               *Dimensions               `json:"dimensions"`
	Cycle                    string                    `json:"cycle"`
	Watering                 string                    `json:"watering"`
	WateringGeneralBenchmark *WateringGeneralBenchmark `json:"watering_general_benchmark"`
	Sunlight                 []string                  `json:"sunlight"`
	PruningMonths            []string                  `json:"pruning_months"`
	PruningCount             *PruningCount             `json:"pruning_count"`
	Soil                     []string                  `json:"soil"`
	Maintenance              string                    `json:"maintenance"`
	PoisonousToHumans        bool                      `json:"poisonous_to_humans"`
	PoisonousToPets          bool                      `json:"poisonous_to_pets"`
	CareLevel                string                    `json:"care_level"`
	PerenualDescription      string                    `json:"perenual_description"`
	Flowers                  bool                      `json:"flowers"`
	FloweringSeason          string                    `json:"flowering_season"`
	WikipediaDescription     *string                   `json:"wikipedia_description"`
	WikipediaExtract         *string                   `json:"wikipedia_extract"`
	WikipediaExtractHTML     *string                   `json:"wikipedia_extract_html"`
	WikipediaURL             *string                   `json:"wikipedia_url"`
}

type Species struct {
	ID              int      `json:"id"`
	CommonName      string   `json:"common_name"`
	ScientificNames []string `json:"scientific_names"`
	OtherNames      []string `json:"other_names"`
	Image           *Image   `json:"image"`
}

type Dimensions struct {
	MinValue float64 `json:"min_value"`
	MaxValue float64 `json:"max_value"`
	Unit     string  `json:"unit"`
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
