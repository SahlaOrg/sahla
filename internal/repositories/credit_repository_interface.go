package repositories

import (
	"context"
	"github.com/mohamed2394/sahla/internal/domains"
)

// CreditApplicationRepository defines the interface for credit application data access
type CreditApplicationRepository interface {
	Create(ctx context.Context, app *domains.CreditApplication) error
	GetByID(ctx context.Context, id uint) (*domains.CreditApplication, error)
	Update(ctx context.Context, app *domains.CreditApplication) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*domains.CreditApplication, int, error)
	GetByUserID(ctx context.Context, userID string) ([]*domains.CreditApplication, error)
}

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	Create(ctx context.Context, payment *domains.Payment) error
	GetByID(ctx context.Context, id uint) (*domains.Payment, error)
	Update(ctx context.Context, payment *domains.Payment) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*domains.Payment, int, error)
	GetByCreditApplicationID(ctx context.Context, creditApplicationID uint) ([]*domains.Payment, error)
	GetByOrderID(ctx context.Context, orderID string) (*domains.Payment, error)
}

// InstallmentRepository defines the interface for installment data access
type InstallmentRepository interface {
	Create(ctx context.Context, installment *domains.Installment) error
	GetByID(ctx context.Context, id uint) (*domains.Installment, error)
	Update(ctx context.Context, installment *domains.Installment) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*domains.Installment, int, error)
	GetByPaymentID(ctx context.Context, paymentID uint) ([]*domains.Installment, error)
}