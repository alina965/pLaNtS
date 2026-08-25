package plant

import (
	"github.com/alina965/pLaNtS/plant-service/internal/client"
	"github.com/alina965/pLaNtS/plant-service/internal/domain"
)

type Service struct {
	perenualClient  *client.PerenualClient
	wikipediaClient *client.WikipediaClient
}

func NewService(perenualClient *client.PerenualClient, wikipediaClient *client.WikipediaClient) *Service {
	return &Service{perenualClient: perenualClient, wikipediaClient: wikipediaClient}
}

func (s *Service) GetSpecies(page int, query string) (*domain.SpeciesList, error) {
	allSpecies, err := s.perenualClient.GetPlants(page, query)
	if err != nil {
		return nil, err
	}

	result := &domain.SpeciesList{PerPage: allSpecies.PerPage, CurrentPage: allSpecies.CurrentPage, LastPage: allSpecies.LastPage, Total: allSpecies.Total}

	for _, item := range allSpecies.Data {
		newSpecies := &domain.Species{ID: item.ID, CommonName: item.CommonName, ScientificNames: item.ScientificName, OtherNames: item.OtherName}

		if item.DefaultImage != nil {
			newSpecies.Image = &domain.Image{
				RegularUrl: item.DefaultImage.RegularUrl,
				SmallUrl:   item.DefaultImage.SmallUrl,
			}
		}

		result.Species = append(result.Species, *newSpecies)
	}

	return result, nil
}

func (s *Service) GetSpeciesDetails(speciesID int) (*domain.SpeciesDetails, error) {
	speciesDetails, err := s.perenualClient.GetPlantDetails(speciesID)
	if err != nil {
		return nil, err
	}

	result := &domain.SpeciesDetails{
		ID:                  speciesID,
		CommonName:          speciesDetails.CommonName,
		ScientificNames:     speciesDetails.ScientificName,
		OtherNames:          speciesDetails.OtherName,
		GrowthRate:          speciesDetails.GrowthRate,
		Type:                speciesDetails.Type,
		Cycle:               speciesDetails.Cycle,
		Watering:            speciesDetails.Watering,
		Sunlight:            speciesDetails.Sunlight,
		PruningMonths:       speciesDetails.PruningMonth,
		Soil:                speciesDetails.Soil,
		Maintenance:         speciesDetails.Maintenance,
		PoisonousToHumans:   speciesDetails.PoisonousToHumans,
		PoisonousToPets:     speciesDetails.PoisonousToPets,
		CareLevel:           speciesDetails.CareLevel,
		PerenualDescription: speciesDetails.Description,
		Flowers:             speciesDetails.Flowers,
		FloweringSeason:     speciesDetails.FloweringSeason,
	}

	if speciesDetails.DefaultImage != nil {
		result.Image = &domain.Image{
			RegularUrl: speciesDetails.DefaultImage.RegularUrl,
			SmallUrl:   speciesDetails.DefaultImage.SmallUrl,
		}
	}

	if speciesDetails.Dimensions.Value != nil {
		result.Dimensions = &domain.Dimensions{
			MinValue: speciesDetails.Dimensions.Value.MinValue,
			MaxValue: speciesDetails.Dimensions.Value.MaxValue,
			Unit:     speciesDetails.Dimensions.Value.Unit,
		}
	}

	if speciesDetails.WateringGeneralBenchmark != nil {
		result.WateringGeneralBenchmark = &domain.WateringGeneralBenchmark{
			Value: speciesDetails.WateringGeneralBenchmark.Value,
			Unit:  speciesDetails.WateringGeneralBenchmark.Unit,
		}
	}

	if speciesDetails.PruningCount != nil {
		result.PruningCount = &domain.PruningCount{
			Amount:   speciesDetails.PruningCount.Amount,
			Interval: speciesDetails.PruningCount.Interval,
		}
	}

	if len(speciesDetails.ScientificName) != 0 {
		speciesDescription := s.wikipediaClient.GetDescription(speciesDetails.ScientificName[0])
		if speciesDescription != nil {
			result.WikipediaDescription = &speciesDescription.Description
			result.WikipediaExtract = &speciesDescription.Extract
			result.WikipediaExtractHTML = &speciesDescription.ExtractHTML
			if speciesDescription.ContentURLs != nil && speciesDescription.ContentURLs.Desktop != nil && len(speciesDescription.ContentURLs.Desktop.Page) != 0 {
				result.WikipediaURL = &speciesDescription.ContentURLs.Desktop.Page
			}
		}
	}

	return result, nil
}
