// internal/service/payment_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/mohamed2394/sahla/modules/credit/domain"
	"github.com/mohamed2394/sahla/modules/credit/dto"
	"github.com/mohamed2394/sahla/modules/credit/repository"
)

type paymentService struct {
	paymentRepo repository.PaymentRepository
	loanRepo    repository.LoanRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, loanRepo repository.LoanRepository) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		loanRepo:    loanRepo,
	}
}

func (s *paymentService) MakePayment(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	loan, err := s.loanRepo.GetLoanByID(ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	if loan.Status != domain.LoanStatusActive {
		return nil, errors.New("loan is not active")
	}

	if req.Amount > loan.RemainingAmount {
		return nil, errors.New("payment amount exceeds remaining loan amount")
	}

	payment := &domain.Payment{
		LoanID:        req.LoanID,
		Amount:        req.Amount,
		PaymentDate:   time.Now(),
		Status:        domain.PaymentStatusCompleted,
		PaymentMethod: req.PaymentMethod,
	}

	if err := s.paymentRepo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	// Update loan
	loan.RemainingAmount -= req.Amount
	loan.TotalPaidAmount += req.Amount
	if loan.RemainingAmount == 0 {
		loan.Status = domain.LoanStatusPaid
	}

	if err := s.loanRepo.UpdateLoan(ctx, loan); err != nil {
		return nil, err
	}

	return convertPaymentToDTO(payment), nil
}

func (s *paymentService) GetPaymentsByLoanID(ctx context.Context, loanID uint) ([]*dto.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetPaymentsByLoanID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	var paymentDTOs []*dto.PaymentResponse
	for _, payment := range payments {
		paymentDTOs = append(paymentDTOs, convertPaymentToDTO(payment))
	}

	return paymentDTOs, nil
}

func convertPaymentToDTO(payment *domain.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:              payment.ID,
		UniversalID:     payment.UniversalID,
		LoanID:          payment.LoanID,
		Amount:          payment.Amount,
		PaymentDate:     payment.PaymentDate,
		Status:          payment.Status,
		PaymentMethod:   payment.PaymentMethod,
		TransactionID:   payment.TransactionID,
		ProcessingFee:   payment.ProcessingFee,
		LateFee:         payment.LateFee,
		EarlyPaymentFee: payment.EarlyPaymentFee,
	}
}
