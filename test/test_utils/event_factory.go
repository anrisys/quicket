package test_utils

import (
	"fmt"
	"time"

	"github.com/anrisys/quicket/internal/event"
	"github.com/anrisys/quicket/pkg/config"
	"github.com/anrisys/quicket/pkg/database"
)

func CreateTestEvent(cfg *config.AppConfig, title string, userID uint) (*event.Event, error) {
	gormDB, err := database.MySQLDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(2 * time.Hour)
	desc := "Test event description"

	event := &event.Event{
		Title: title,
		Description: &desc,
		StartDate: startDate,
		EndDate: endDate,
		MaxSeats: 100,
		OrganizerID: userID,
	}

	if err := gormDB.Create(event).Error; err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func CleanupTestEvent(cfg *config.AppConfig, title string) error {
	gormDB, err := database.MySQLDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize db: %w", err)
	}

	return gormDB.Where("title = ?", title).Delete(&event.Event{}).Error
}