package db

import (
	"github.com/mohamed2394/sahla/modules/user/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect initializes a connection to the PostgreSQL database
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}

// AutoMigrateModels performs automatic migration for all models
func AutoMigrateModels() error {
	err := DB.AutoMigrate(&domain.User{})
	if err != nil {
		return err
	}
	return nil
}
