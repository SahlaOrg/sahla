package repository

import (
	"context"

	"github.com/mohamed2394/sahla/modules/credit/domain"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *domain.Loan) error
	GetLoanByID(ctx context.Context, id uint) (*domain.Loan, error)
	GetLoanByUniversalID(ctx context.Context, universalID string) (*domain.Loan, error)
	UpdateLoan(ctx context.Context, loan *domain.Loan) error
	GetLoansByUserID(ctx context.Context, userID uint) ([]*domain.Loan, error)
	GetAllLoans(ctx context.Context, limit, offset int) ([]*domain.Loan, error)
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	GetPaymentByID(ctx context.Context, id uint) (*domain.Payment, error)
	GetPaymentsByLoanID(ctx context.Context, loanID uint) ([]*domain.Payment, error)
	UpdatePayment(ctx context.Context, payment *domain.Payment) error
}

type RepaymentPlanRepository interface {
	CreateRepaymentPlan(ctx context.Context, repaymentPlan *domain.RepaymentPlan) error
	GetRepaymentPlansByLoanID(ctx context.Context, loanID uint) ([]*domain.RepaymentPlan, error)
	UpdateRepaymentPlan(ctx context.Context, repaymentPlan *domain.RepaymentPlan) error
}
