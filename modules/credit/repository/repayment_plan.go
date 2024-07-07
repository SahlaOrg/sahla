package repository

import (
	"context"

	"github.com/mohamed2394/sahla/modules/credit/domain"

	"gorm.io/gorm"
)

type repaymentPlanRepository struct {
	db *gorm.DB
}

func NewRepaymentPlanRepository(db *gorm.DB) RepaymentPlanRepository {
	return &repaymentPlanRepository{db: db}
}

func (r *repaymentPlanRepository) CreateRepaymentPlan(ctx context.Context, repaymentPlan *domain.RepaymentPlan) error {
	return r.db.WithContext(ctx).Create(repaymentPlan).Error
}

func (r *repaymentPlanRepository) GetRepaymentPlansByLoanID(ctx context.Context, loanID uint) ([]*domain.RepaymentPlan, error) {
	var repaymentPlans []*domain.RepaymentPlan
	err := r.db.WithContext(ctx).Where("loan_id = ?", loanID).Order("installment_num").Find(&repaymentPlans).Error
	return repaymentPlans, err
}

func (r *repaymentPlanRepository) UpdateRepaymentPlan(ctx context.Context, repaymentPlan *domain.RepaymentPlan) error {
	return r.db.WithContext(ctx).Save(repaymentPlan).Error
}
