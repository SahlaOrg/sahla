package service

import (
	"context"

	"github.com/mohamed2394/sahla/modules/credit/dto"
)

type LoanService interface {
	CreateLoan(ctx context.Context, req dto.CreateLoanRequest) (*dto.LoanResponse, error)
	GetLoanByID(ctx context.Context, id uint) (*dto.LoanDetailResponse, error)
	GetLoanByUniversalID(ctx context.Context, universalID string) (*dto.LoanDetailResponse, error)
	GetLoansByUserID(ctx context.Context, userID uint) ([]*dto.LoanResponse, error)
	UpdateLoanStatus(ctx context.Context, id uint, status string) error
}

type PaymentService interface {
	MakePayment(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	GetPaymentsByLoanID(ctx context.Context, loanID uint) ([]*dto.PaymentResponse, error)
}

type RepaymentPlanService interface {
	GenerateRepaymentPlan(ctx context.Context, loanID uint) error
	GetRepaymentPlanByLoanID(ctx context.Context, loanID uint) ([]*dto.RepaymentPlanResponse, error)
}
