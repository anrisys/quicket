package test_utils

import (
	"context"
	"fmt"

	"github.com/anrisys/quicket/internal/user"
	"github.com/anrisys/quicket/pkg/config"
	"github.com/anrisys/quicket/pkg/database"
	"github.com/anrisys/quicket/pkg/security"
)

func CreateTestUser(cfg *config.AppConfig) (*user.User, error) {
	accSecurity := security.NewAccountSecurity(cfg)
	hashedPassword, err := accSecurity.HashPassword(context.Background(), "Pass345!@#")
	if err != nil {
		return nil, fmt.Errorf("CreateTestUser: failed to hash user's example pass: %w", err)
	}
	publicID, err := accSecurity.GeneratePublicID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("CreateTestUser: failed to generate public id: %w", err)
	}
	data := &user.User{
		PublicID: publicID,
		Email: "user@example.com",
		Password: hashedPassword,
		Role: "admin",
	}
	gormDB,err := database.MySQLDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db %w", err)
	}
	if err := gormDB.Create(data).Error; err != nil {
        return nil, fmt.Errorf("failed to seed user: %w", err)
    }
    return data, nil
}

func CleanupTestUser(cfg *config.AppConfig, email string) error {
	db, err := database.MySQLDB(cfg)
	if err != nil {
		return err
	}
	
	return db.Unscoped().Where("email = ?", email).Delete(&user.User{}).Error
}