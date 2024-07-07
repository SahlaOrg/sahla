// internal/repository/gorm/payment_repository.go
package repository

import (
	"context"

	"github.com/mohamed2394/sahla/modules/credit/domain"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetPaymentByID(ctx context.Context, id uint) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetPaymentsByLoanID(ctx context.Context, loanID uint) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	err := r.db.WithContext(ctx).Where("loan_id = ?", loanID).Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}
