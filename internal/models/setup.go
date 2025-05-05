package models

import (
	"github.com/kartikey1188/build-in-progress_01/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupGorm(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.StoragePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}
