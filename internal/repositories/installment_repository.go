package repositories

import (
	"context"
	"errors"

	"github.com/mohamed2394/sahla/internal/domains"
	"gorm.io/gorm"
)

type installmentRepository struct {
	db *gorm.DB
}

func NewInstallmentRepository(db *gorm.DB) InstallmentRepository {
	return &installmentRepository{db: db}
}

func (r *installmentRepository) Create(ctx context.Context, installment *domains.Installment) error {
	return r.db.WithContext(ctx).Create(installment).Error
}

func (r *installmentRepository) GetByID(ctx context.Context, id uint) (*domains.Installment, error) {
	var installment domains.Installment
	if err := r.db.WithContext(ctx).First(&installment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("installment not found")
		}
		return nil, err
	}
	return &installment, nil
}

func (r *installmentRepository) Update(ctx context.Context, installment *domains.Installment) error {
	return r.db.WithContext(ctx).Save(installment).Error
}

func (r *installmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domains.Installment{}, id).Error
}

func (r *installmentRepository) List(ctx context.Context, offset, limit int) ([]*domains.Installment, int, error) {
	var installments []*domains.Installment
	var total int64

	err := r.db.WithContext(ctx).Model(&domains.Installment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&installments).Error
	if err != nil {
		return nil, 0, err
	}

	return installments, int(total), nil
}

func (r *installmentRepository) GetByPaymentID(ctx context.Context, paymentID uint) ([]*domains.Installment, error) {
	var installments []*domains.Installment
	err := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).Find(&installments).Error
	return installments, err
}