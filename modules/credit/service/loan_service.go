// internal/service/loan_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/mohamed2394/sahla/modules/credit/domain"
	"github.com/mohamed2394/sahla/modules/credit/dto"
	"github.com/mohamed2394/sahla/modules/credit/repository"
)

type loanService struct {
	loanRepo          repository.LoanRepository
	repaymentPlanRepo repository.RepaymentPlanRepository
}

func NewLoanService(loanRepo repository.LoanRepository, repaymentPlanRepo repository.RepaymentPlanRepository) LoanService {
	return &loanService{
		loanRepo:          loanRepo,
		repaymentPlanRepo: repaymentPlanRepo,
	}
}

func (s *loanService) CreateLoan(ctx context.Context, req dto.CreateLoanRequest) (*dto.LoanResponse, error) {
	// TODO: Implement credit score check and approval logic here

	loan := &domain.Loan{
		Amount:          req.Amount,
		OriginalTerm:    req.Term,
		RemainingTerm:   req.Term,
		Status:          domain.LoanStatusPending,
		ApprovedAt:      nil, // Will be set when approved
		ActivatedAt:     nil, // Will be set when activated
		RemainingAmount: req.Amount,
	}

	if err := s.loanRepo.CreateLoan(ctx, loan); err != nil {
		return nil, err
	}

	// Generate repayment plan
	if err := s.generateRepaymentPlan(ctx, loan); err != nil {
		return nil, err
	}

	return convertLoanToDTO(loan), nil
}

func (s *loanService) GetLoanByID(ctx context.Context, id uint) (*dto.LoanDetailResponse, error) {
	loan, err := s.loanRepo.GetLoanByID(ctx, id)
	if err != nil {
		return nil, err
	}

	repaymentPlans, err := s.repaymentPlanRepo.GetRepaymentPlansByLoanID(ctx, loan.ID)
	if err != nil {
		return nil, err
	}

	return convertLoanToDetailDTO(loan, repaymentPlans), nil
}

func (s *loanService) GetLoanByUniversalID(ctx context.Context, universalID string) (*dto.LoanDetailResponse, error) {
	loan, err := s.loanRepo.GetLoanByUniversalID(ctx, universalID)
	if err != nil {
		return nil, err
	}

	repaymentPlans, err := s.repaymentPlanRepo.GetRepaymentPlansByLoanID(ctx, loan.ID)
	if err != nil {
		return nil, err
	}

	return convertLoanToDetailDTO(loan, repaymentPlans), nil
}

func (s *loanService) GetLoansByUserID(ctx context.Context, userID uint) ([]*dto.LoanResponse, error) {
	loans, err := s.loanRepo.GetLoansByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var loanDTOs []*dto.LoanResponse
	for _, loan := range loans {
		loanDTOs = append(loanDTOs, convertLoanToDTO(loan))
	}

	return loanDTOs, nil
}

func (s *loanService) UpdateLoanStatus(ctx context.Context, id uint, status string) error {
	loan, err := s.loanRepo.GetLoanByID(ctx, id)
	if err != nil {
		return err
	}

	newStatus := domain.LoanStatus(status)
	if !isValidLoanStatus(newStatus) {
		return errors.New("invalid loan status")
	}

	loan.Status = newStatus

	if newStatus == domain.LoanStatusApproved {
		now := time.Now()
		loan.ApprovedAt = &now
	} else if newStatus == domain.LoanStatusActive {
		now := time.Now()
		loan.ActivatedAt = &now
	}

	return s.loanRepo.UpdateLoan(ctx, loan)
}

func (s *loanService) generateRepaymentPlan(ctx context.Context, loan *domain.Loan) error {
	// Calculate fixed monthly payment
	monthlyPayment := loan.Amount / float64(loan.OriginalTerm)
	remainingPrincipal := loan.Amount

	for i := 1; i <= loan.OriginalTerm; i++ {
		principalPayment := monthlyPayment
		if i == loan.OriginalTerm {
			// Adjust the last payment to account for any rounding errors
			principalPayment = remainingPrincipal
		}

		repaymentPlan := &domain.RepaymentPlan{
			LoanID:         loan.ID,
			InstallmentNum: i,
			DueDate:        loan.ActivatedAt.AddDate(0, i, 0),
			Amount:         principalPayment,
			PrincipalPart:  principalPayment,
			Status:         "PENDING",
		}

		if err := s.repaymentPlanRepo.CreateRepaymentPlan(ctx, repaymentPlan); err != nil {
			return err
		}

		remainingPrincipal -= principalPayment
	}

	return nil
}

func isValidLoanStatus(status domain.LoanStatus) bool {
	validStatuses := []domain.LoanStatus{
		domain.LoanStatusPending,
		domain.LoanStatusApproved,
		domain.LoanStatusRejected,
		domain.LoanStatusActive,
		domain.LoanStatusPaid,
		domain.LoanStatusDefaulted,
	}

	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

func convertLoanToDTO(loan *domain.Loan) *dto.LoanResponse {
	return &dto.LoanResponse{
		ID:                loan.ID,
		UniversalID:       loan.UniversalID,
		Amount:            loan.Amount,
		OriginalTerm:      loan.OriginalTerm,
		RemainingTerm:     loan.RemainingTerm,
		Status:            loan.Status,
		ApprovedAt:        loan.ApprovedAt,
		ActivatedAt:       loan.ActivatedAt,
		NextPaymentDate:   loan.NextPaymentDate,
		RemainingAmount:   loan.RemainingAmount,
		TotalPaidAmount:   loan.TotalPaidAmount,
		CreditScoreImpact: loan.CreditScoreImpact,
	}
}

func convertLoanToDetailDTO(loan *domain.Loan, repaymentPlans []*domain.RepaymentPlan) *dto.LoanDetailResponse {
	loanDTO := convertLoanToDTO(loan)
	repaymentPlanDTOs := make([]dto.RepaymentPlanResponse, len(repaymentPlans))

	for i, plan := range repaymentPlans {
		repaymentPlanDTOs[i] = dto.RepaymentPlanResponse{
			ID:             plan.ID,
			LoanID:         plan.LoanID,
			InstallmentNum: plan.InstallmentNum,
			DueDate:        plan.DueDate,
			Amount:         plan.Amount,
			PrincipalPart:  plan.PrincipalPart,
			Status:         plan.Status,
		}
	}

	return &dto.LoanDetailResponse{
		LoanResponse:      *loanDTO,
		RepaymentSchedule: repaymentPlanDTOs,
	}
}
