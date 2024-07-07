// internal/service/repayment_plan_service.go
package service

import (
	"context"

	"github.com/mohamed2394/sahla/modules/credit/dto"
	"github.com/mohamed2394/sahla/modules/credit/repository"
)

type repaymentPlanService struct {
	repaymentPlanRepo repository.RepaymentPlanRepository
	loanRepo          repository.LoanRepository
}

func NewRepaymentPlanService(repaymentPlanRepo repository.RepaymentPlanRepository, loanRepo repository.LoanRepository) RepaymentPlanService {
	return &repaymentPlanService{
		repaymentPlanRepo: repaymentPlanRepo,
		loanRepo:          loanRepo,
	}
}

func (s *repaymentPlanService) GenerateRepaymentPlan(ctx context.Context, loanID uint) error {
	// This method is called from the LoanService when a loan is created
	// The implementation is in the LoanService to avoid circular dependencies
	return nil
}

func (s *repaymentPlanService) GetRepaymentPlanByLoanID(ctx context.Context, loanID uint) ([]*dto.RepaymentPlanResponse, error) {
	repaymentPlans, err := s.repaymentPlanRepo.GetRepaymentPlansByLoanID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	var repaymentPlanDTOs []*dto.RepaymentPlanResponse
	for _, plan := range repaymentPlans {
		repaymentPlanDTOs = append(repaymentPlanDTOs, &dto.RepaymentPlanResponse{
			ID:             plan.ID,
			LoanID:         plan.LoanID,
			InstallmentNum: plan.InstallmentNum,
			DueDate:        plan.DueDate,
			Amount:         plan.Amount,
			PrincipalPart:  plan.PrincipalPart,
			Status:         plan.Status,
		})
	}

	return repaymentPlanDTOs, nil
}
