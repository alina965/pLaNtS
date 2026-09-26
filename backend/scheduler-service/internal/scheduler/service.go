package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/alina965/pLaNtS/scheduler-service/internal/client"
	"github.com/robfig/cron/v3"
)

type Service struct {
	location string
	client   *client.UserPlantsClient
}

func NewService(location string, client *client.UserPlantsClient) *Service {
	return &Service{location: location, client: client}
}

func (s *Service) Setup() (*cron.Cron, error) {
	location, err := time.LoadLocation(s.location)
	if err != nil {
		return nil, err
	}
	c := cron.New(cron.WithLocation(location))
	_, err = c.AddFunc("0 9 * * *", s.processDuePlants)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Service) processDuePlants() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	plants, err := s.client.ListNeedingWater(ctx)
	if err != nil {
		log.Printf("list needing water: %v", err)
		return
	}

	for _, plant := range plants {
		var newStatus string
		switch plant.Status {
		case "ok":
			newStatus = "due"
		case "due":
			newStatus = "overdue"
		default:
			continue
		}
		if err := s.client.UpdateStatus(ctx, plant.ID, newStatus); err != nil {
			log.Printf("update status %s: %v", plant.ID, err)
			continue
		}
		// TODO: telegram notify
	}
}
