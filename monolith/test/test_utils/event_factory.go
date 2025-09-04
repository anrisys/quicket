package test_utils

import (
	"context"
	"fmt"
	"time"

	"github.com/anrisys/quicket/internal/event"
	"github.com/anrisys/quicket/pkg/config"
	"github.com/anrisys/quicket/pkg/database"
	"github.com/anrisys/quicket/pkg/security"
)

func CreateTestEvent(cfg *config.AppConfig, title string, userID uint) (*event.Event, error) {
	gormDB, err := database.MySQLDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(2 * time.Hour)
	desc := "Test event description"

	as := security.NewAccountSecurity(cfg)
	publicID, err := as.GeneratePublicID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to generate public id for test event: %w", err)
	}

	event := &event.Event{
		Title: title,
		PublicID: publicID,
		Description: &desc,
		StartDate: startDate,
		EndDate: endDate,
		MaxSeats: 100,
		AvailableSeats: 100,
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

func CreatePastEventTest(cfg *config.AppConfig, title string, userID uint) (*event.Event, error) {
	gormDB, err := database.MySQLDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	startDate := time.Now().AddDate(0, 0, -5)
	endDate := time.Now().AddDate(0, 0, -4)

	as := security.NewAccountSecurity(cfg)
	publicID, err := as.GeneratePublicID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to generate public id for test event: %w", err)
	}

	desc := "A past event"

	pastEvent := &event.Event{
		PublicID: publicID,
		Title: title,
		Description: &desc,
		StartDate: startDate,
		EndDate: endDate,
		MaxSeats: 100,
		AvailableSeats: 100,
		OrganizerID: userID,
	}

	if err := gormDB.Create(pastEvent).Error; err != nil {
		return nil, fmt.Errorf("failed to create event test: %w", err)
	}

	return pastEvent, nil
}

func CreateLimitedSeatsEventTest(cfg *config.AppConfig, title string, userID uint) (*event.Event, error) {
	gormDB, err := database.MySQLDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	startDate := time.Now().AddDate(0, 0, 10)
	endDate := time.Now().AddDate(0, 0, 11)

	as := security.NewAccountSecurity(cfg)
	publicID, err := as.GeneratePublicID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to generate public id for test event: %w", err)
	}

	desc := "Event with limited seats"

	pastEvent := &event.Event{
		PublicID: publicID,
		Title: title,
		Description: &desc,
		StartDate: startDate,
		EndDate: endDate,
		MaxSeats: 5,
		AvailableSeats: 5,
		OrganizerID: userID,
	}

	if err := gormDB.Create(pastEvent).Error; err != nil {
		return nil, fmt.Errorf("failed to create event test: %w", err)
	}

	return pastEvent, nil
}