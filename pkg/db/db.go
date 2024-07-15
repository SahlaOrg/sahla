package db

import (
	"log"
	domain "github.com/mohamed2394/sahla/internal/domains"
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
	err := dbInstance.AutoMigrate(
		&domain.User{},
		&domain.CreditApplication{},
		&domain.Payment{},
		&domain.Installment{},
	)
	if err != nil {
		return err
	}

	// Ensure indexes are created for foreign keys and unique constraints
	err = dbInstance.Exec("CREATE INDEX IF NOT EXISTS idx_payments_credit_application_id ON payments(credit_application_id)").Error
	if err != nil {
		return err
	}

	err = dbInstance.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id)").Error
	if err != nil {
		return err
	}

	err = dbInstance.Exec("CREATE INDEX IF NOT EXISTS idx_installments_payment_id ON installments(payment_id)").Error
	if err != nil {
		return err
	}

	return nil
}

// GetDB returns the instance of the database connection
func GetDB() *gorm.DB {
	if dbInstance == nil {
		log.Fatal("Database connection is not initialized")
	}
	return dbInstance
}