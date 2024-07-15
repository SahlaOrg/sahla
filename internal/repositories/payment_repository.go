package repositories
import (
		"gorm.io/gorm"
"context"
	"github.com/mohamed2394/sahla/internal/domains"
"errors"


)
type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domains.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id uint) (*domains.Payment, error) {
	var payment domains.Payment
	if err := r.db.WithContext(ctx).First(&payment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *domains.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domains.Payment{}, id).Error
}

func (r *paymentRepository) List(ctx context.Context, offset, limit int) ([]*domains.Payment, int, error) {
	var payments []*domains.Payment
	var total int64

	err := r.db.WithContext(ctx).Model(&domains.Payment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, int(total), nil
}

func (r *paymentRepository) GetByCreditApplicationID(ctx context.Context, creditApplicationID uint) ([]*domains.Payment, error) {
	var payments []*domains.Payment
	err := r.db.WithContext(ctx).Where("credit_application_id = ?", creditApplicationID).Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domains.Payment, error) {
	var payment domains.Payment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found for the given order ID")
		}
		return nil, err
	}
	return &payment, nil
}
