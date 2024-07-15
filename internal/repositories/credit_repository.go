package repositories

import (
	"context"
	"errors"
	utils "github.com/mohamed2394/sahla/internal/utils"

	"github.com/mohamed2394/sahla/internal/domains"
	"gorm.io/gorm"
)

type creditApplicationRepository struct {
	db *gorm.DB
}

// NewCreditApplicationRepository creates a new instance of CreditApplicationRepository
func NewCreditApplicationRepository(db *gorm.DB) CreditApplicationRepository {
	return &creditApplicationRepository{db: db}
}

func (r *creditApplicationRepository) Create(ctx context.Context, app *domains.CreditApplication) error {
	err := r.db.WithContext(ctx).Create(app).Error
	if err != nil {
		return &utils.ErrDatabase{Err: err}
	}
	return nil
}

func (r *creditApplicationRepository) GetByID(ctx context.Context, id uint) (*domains.CreditApplication, error) {
	var app domains.CreditApplication
	if err := r.db.WithContext(ctx).First(&app, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utils.ErrNotFound{Entity: "CreditApplication", ID: id}
		}
		return nil, &utils.ErrDatabase{Err: err}
	}
	return &app, nil
}

func (r *creditApplicationRepository) Update(ctx context.Context, app *domains.CreditApplication) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *creditApplicationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domains.CreditApplication{}, id).Error
}

func (r *creditApplicationRepository) List(ctx context.Context, offset, limit int) ([]*domains.CreditApplication, int, error) {
	var apps []*domains.CreditApplication
	var total int64

	err := r.db.WithContext(ctx).Model(&domains.CreditApplication{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, int(total), nil
}

func (r *creditApplicationRepository) GetByUserID(ctx context.Context, userID string) ([]*domains.CreditApplication, error) {
	var apps []*domains.CreditApplication
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&apps).Error
	return apps, err
}