package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alina965/pLaNtS/scheduler-service/internal/client"
	"github.com/alina965/pLaNtS/scheduler-service/internal/domain"
	"github.com/alina965/pLaNtS/scheduler-service/internal/kafka"
	"github.com/robfig/cron/v3"
)

type Service struct {
	location string
	plants   *client.UserPlantsClient
	producer *kafka.WateringProducer
}

func NewService(location string, plants *client.UserPlantsClient, producer *kafka.WateringProducer) *Service {
	return &Service{location: location, plants: plants, producer: producer}
}

func (s *Service) Setup() (*cron.Cron, error) {
	location, err := time.LoadLocation(s.location)
	if err != nil {
		return nil, err
	}
	c := cron.New(cron.WithLocation(location))
	_, err = c.AddFunc("* * * * *", s.processDuePlants)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Service) processDuePlants() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	plants, err := s.plants.ListNeedingWater(ctx)
	if err != nil {
		log.Printf("list needing water: %v", err)
		return
	}

	for _, plant := range plants {
		var newStatus string
		var message string
		switch plant.Status {
		case "ok":
			newStatus = "due"
			message = fmt.Sprintf("Пора полить растение «%s»", plant.Name)
		case "due":
			newStatus = "overdue"
			message = fmt.Sprintf("Полив просрочен: «%s»", plant.Name)
		default:
			continue
		}

		if err := s.plants.UpdateStatus(ctx, plant.ID, newStatus); err != nil {
			log.Printf("update status %s: %v", plant.ID, err)
			continue
		}

		event := domain.WateringNotifyEvent{
			UserID:     plant.UserID,
			PlantID:    plant.ID,
			PlantName:  plant.Name,
			Status:     newStatus,
			Text:       message,
			OccurredAt: time.Now().UTC(),
		}

		if err := s.producer.Publish(ctx, event); err != nil {
			log.Printf("publish watering notify: %v", err)
			continue
		}
	}
}
