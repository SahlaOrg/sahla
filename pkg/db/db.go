package db

import (
	"log"

	"github.com/mohamed2394/sahla/modules/user/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

// Connect initializes the database connection
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	dbInstance = db

	return db, nil
}

// AutoMigrateModels migrates the database models
func AutoMigrateModels() error {
	return dbInstance.AutoMigrate(&domain.User{})
}

// GetDB returns the instance of the database connection
func GetDB() *gorm.DB {
	if dbInstance == nil {
		log.Fatal("Database connection is not initialized")
	}
	return dbInstance
}
