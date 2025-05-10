package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kartikey1188/build-in-progress_01/internal/config"
	"github.com/kartikey1188/build-in-progress_01/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Postgres struct {
	GormDB *gorm.DB
	SqlDB  *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	gormDB, err := models.SetupGorm(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get *sql.DB from GORM: %w", err)
	}

	if err := autoMigrateTables(gormDB); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate tables: %w", err)
	}

	if err := createAdminUser(gormDB); err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	return &Postgres{
		GormDB: gormDB,
		SqlDB:  sqlDB,
	}, nil
}

func autoMigrateTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Business{},
		&models.Collector{},
		&models.ServiceCategory{},
		&models.CollectorServiceCategory{},
		&models.Vehicle{},
		&models.CollectorDriver{},
		&models.CollectorVehicle{},
		&models.VehicleDriver{},
		&models.PickupRequest{},
	)
}

func createAdminUser(db *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	adminUser := models.User{
		Email:        "admin@gmail.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Application Admin",
		Role:         "Admin",
		IsActive:     true,
		IsVerified:   true,
		IsFlagged:    false,
		PhoneNumber:  "10000000",
		Address:      "N/A",
		ProfileImage: "default.jpg",
		Registration: time.Now(),
		LastLogin:    time.Now(),
	}

	result := db.Where(models.User{Email: adminUser.Email}).FirstOrCreate(&adminUser)
	if result.Error != nil {
		return fmt.Errorf("failed to create admin user: %w", result.Error)
	}

	return nil
}
