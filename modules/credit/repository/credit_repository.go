// internal/repository/gorm/loan_repository.go
package repository

import (
	"context"

	"github.com/mohamed2394/sahla/modules/credit/domain"
	"gorm.io/gorm"
)

type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) CreateLoan(ctx context.Context, loan *domain.Loan) error {
	return r.db.WithContext(ctx).Create(loan).Error
}

func (r *loanRepository) GetLoanByID(ctx context.Context, id uint) (*domain.Loan, error) {
	var loan domain.Loan
	err := r.db.WithContext(ctx).First(&loan, id).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) GetLoanByUniversalID(ctx context.Context, universalID string) (*domain.Loan, error) {
	var loan domain.Loan
	err := r.db.WithContext(ctx).Where("universal_id = ?", universalID).First(&loan).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) UpdateLoan(ctx context.Context, loan *domain.Loan) error {
	return r.db.WithContext(ctx).Save(loan).Error
}

func (r *loanRepository) GetLoansByUserID(ctx context.Context, userID uint) ([]*domain.Loan, error) {
	var loans []*domain.Loan
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&loans).Error
	return loans, err
}

func (r *loanRepository) GetAllLoans(ctx context.Context, limit, offset int) ([]*domain.Loan, error) {
	var loans []*domain.Loan
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&loans).Error
	return loans, err
}
