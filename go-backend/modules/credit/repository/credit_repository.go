package repositories

import (
	credit "github.com/mohamed2394/sahla/modules/credit"
	"gorm.io/gorm"
)

type GormCreditRepository struct {
	db *gorm.DB
}

func NewGormCreditRepository(db *gorm.DB) *GormCreditRepository {
	return &GormCreditRepository{db: db}
}

func (r *GormCreditRepository) SaveCreditScore(score *credit.CreditScore) error {
	return r.db.Save(score).Error
}

func (r *GormCreditRepository) GetCreditScoreByUserID(userID string) (*credit.CreditScore, error) {
	var score credit.CreditScore
	err := r.db.Where("user_id = ?", userID).First(&score).Error
	return &score, err
}

func (r *GormCreditRepository) SaveCreditFeatures(features *credit.CreditFeatures) error {
	return r.db.Save(features).Error
}

func (r *GormCreditRepository) GetCreditFeaturesByUserID(userID string) (*credit.CreditFeatures, error) {
	var features credit.CreditFeatures
	err := r.db.Where("user_id = ?", userID).First(&features).Error
	return &features, err
}

// Implement other methods for GormCreditRepository...
