package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	repository "github.com/mohamed2394/sahla/internal/repositories"

	"github.com/mohamed2394/sahla/internal/domains"
	"github.com/mohamed2394/sahla/internal/utils"
	"go.uber.org/zap"
)

var (
	ErrInsufficientCredit = errors.New("insufficient credit")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrPaymentFailed      = errors.New("payment processing failed")
)
type CreditPaymentServiceInterface interface {
	CreateCreditApplication(ctx context.Context, app *domains.CreditApplication) error
	ApproveCreditApplication(ctx context.Context, id uint) error
	CreatePayment(ctx context.Context, payment *domains.Payment) error
	HandlePaymentWebhook(ctx context.Context, paymentID uint, status string) error
	ProcessPayment(ctx context.Context, paymentID uint) error
	ProcessInstallment(ctx context.Context, installmentID uint) error
	HandleInstallmentWebhook(ctx context.Context, installmentID uint, status string) error
	GetPaymentDetails(ctx context.Context, paymentID uint) (*domains.Payment, error)
}

type CreditPaymentService struct {
	creditAppRepo   repository.CreditApplicationRepository
	paymentRepo     repository.PaymentRepository
	installmentRepo repository.InstallmentRepository
	logger          *zap.Logger
	paymentGateway  PaymentGateway
}

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, amount int, currency string, paymentMethod domains.PaymentMethod) error
	SimulatePaymentWebhook(ctx context.Context, paymentID uint) error
	SimulateInstallmentWebhook(ctx context.Context, installmentID uint) error
}

func NewCreditPaymentService(
	creditAppRepo repository.CreditApplicationRepository,
	paymentRepo repository.PaymentRepository,
	installmentRepo repository.InstallmentRepository,
	logger *zap.Logger,
	paymentGateway PaymentGateway,
) *CreditPaymentService {
	return &CreditPaymentService{
		creditAppRepo:   creditAppRepo,
		paymentRepo:     paymentRepo,
		installmentRepo: installmentRepo,
		logger:          logger,
		paymentGateway:  paymentGateway,
	}
}

func (s *CreditPaymentService) CreateCreditApplication(ctx context.Context, app *domains.CreditApplication) error {
	s.logger.Info("Creating credit application", zap.Any("application", app))
	
	if app.Amount <= 0 {
		return ErrInvalidAmount
	}
	
	app.Status = "PENDING"
	err := s.creditAppRepo.Create(ctx, app)
	if err != nil {
		s.logger.Error("Failed to create credit application", zap.Error(err))
		return fmt.Errorf("failed to create credit application: %w", err)
	}
	
	s.logger.Info("Credit application created successfully", zap.Uint("id", app.ID))
	return nil
}

func (s *CreditPaymentService) ApproveCreditApplication(ctx context.Context, id uint) error {
	s.logger.Info("Approving credit application", zap.Uint("id", id))
	
	app, err := s.creditAppRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get credit application", zap.Error(err))
		return fmt.Errorf("failed to get credit application: %w", err)
	}
	
	// Perform credit check (this is a simplified example)
	creditScore, err := s.performCreditCheck(ctx, app.UserID)
	if err != nil {
		s.logger.Error("Failed to perform credit check", zap.Error(err))
		return fmt.Errorf("failed to perform credit check: %w", err)
	}
	
	if creditScore < 650 {
		app.Status = "REJECTED"
		s.logger.Info("Credit application rejected due to low credit score", zap.Int("creditScore", creditScore))
	} else {
		app.Status = "APPROVED"
		s.logger.Info("Credit application approved", zap.Int("creditScore", creditScore))
	}
	
	err = s.creditAppRepo.Update(ctx, app)
	if err != nil {
		s.logger.Error("Failed to update credit application", zap.Error(err))
		return fmt.Errorf("failed to update credit application: %w", err)
	}
	
	return nil
}

func (s *CreditPaymentService) performCreditCheck(ctx context.Context, userID string) (int, error) {
    // In a real-world scenario, this would involve calling a credit bureau API
    // For this example, we'll use our weighted random credit score generator
    creditScore := utils.GenerateWeightedRandomCreditScore()
    
    s.logger.Info("Performed credit check", 
        zap.String("userID", userID), 
        zap.Int("creditScore", creditScore))
    
    return creditScore, nil
}
func (s *CreditPaymentService) CreatePayment(ctx context.Context, payment *domains.Payment) error {
	s.logger.Info("Creating payment", zap.Any("payment", payment))
	
	if payment.Amount <= 0 {
		return ErrInvalidAmount
	}
	
	// Check if the user has sufficient credit
	creditApp, err := s.creditAppRepo.GetByID(ctx, payment.CreditApplicationID)
	if err != nil {
		s.logger.Error("Failed to get credit application", zap.Error(err))
		return fmt.Errorf("failed to get credit application: %w", err)
	}
	
	if creditApp.Status != "APPROVED" {
		s.logger.Warn("Attempt to create payment for unapproved credit application", zap.Uint("creditAppID", creditApp.ID))
		return errors.New("credit application not approved")
	}
	
	totalPayments, err := s.getTotalPaymentsForCreditApp(ctx, payment.CreditApplicationID)
	if err != nil {
		s.logger.Error("Failed to get total payments", zap.Error(err))
		return fmt.Errorf("failed to get total payments: %w", err)
	}
	
	if totalPayments+payment.Amount > creditApp.Amount {
		s.logger.Warn("Insufficient credit for payment", 
			zap.Int("requestedAmount", payment.Amount), 
			zap.Int("availableCredit", creditApp.Amount-totalPayments))
		return ErrInsufficientCredit
	}
	
	payment.Status = "PENDING"
	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		s.logger.Error("Failed to create payment", zap.Error(err))
		return fmt.Errorf("failed to create payment: %w", err)
	}
	
	s.logger.Info("Payment created successfully", zap.Uint("id", payment.ID))
	go func() {
		time.Sleep(5 * time.Second) // Simulate a delay
		err := s.paymentGateway.SimulatePaymentWebhook(context.Background(), payment.ID)
		if err != nil {
			s.logger.Error("Failed to simulate payment webhook", zap.Error(err))
		}
	}()

	return nil

}
func (s *CreditPaymentService) HandlePaymentWebhook(ctx context.Context, paymentID uint, status string) error {
	s.logger.Info("Handling payment webhook", zap.Uint("paymentID", paymentID), zap.String("status", status))
	
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		s.logger.Error("Failed to get payment", zap.Error(err))
		return fmt.Errorf("failed to get payment: %w", err)
	}
	
	if status == "SUCCESSFUL" {
		payment.Status = "SUCCESSFUL"
		
		installments, err := s.createInstallments(ctx, payment)
		if err != nil {
			s.logger.Error("Failed to create installments", zap.Error(err))
			return fmt.Errorf("failed to create installments: %w", err)
		}
		
		payment.Installments = installments
	} else if status == "FAILED" {
		payment.Status = "FAILED"
	} else {
		return fmt.Errorf("unknown payment status: %s", status)
	}
	
	err = s.paymentRepo.Update(ctx, payment)
	if err != nil {
		s.logger.Error("Failed to update payment", zap.Error(err))
		return fmt.Errorf("failed to update payment: %w", err)
	}
	
	s.logger.Info("Payment webhook handled successfully", zap.Uint("paymentID", paymentID), zap.String("status", status))
	return nil
}


func (s *CreditPaymentService) ProcessPayment(ctx context.Context, paymentID uint) error {
	s.logger.Info("Processing payment", zap.Uint("paymentID", paymentID))
	
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		s.logger.Error("Failed to get payment", zap.Error(err))
		return fmt.Errorf("failed to get payment: %w", err)
	}
	
	err = s.paymentGateway.ProcessPayment(ctx, payment.Amount, payment.Currency, payment.PaymentMethod)
	if err != nil {
		s.logger.Error("Payment processing failed", zap.Error(err))
		payment.Status = "FAILED"
		_ = s.paymentRepo.Update(ctx, payment)
		return ErrPaymentFailed
	}
	
	payment.Status = "SUCCESSFUL"
	
	installments, err := s.createInstallments(ctx, payment)
	if err != nil {
		s.logger.Error("Failed to create installments", zap.Error(err))
		return fmt.Errorf("failed to create installments: %w", err)
	}
	
	payment.Installments = installments
	err = s.paymentRepo.Update(ctx, payment)
	if err != nil {
		s.logger.Error("Failed to update payment", zap.Error(err))
		return fmt.Errorf("failed to update payment: %w", err)
	}
	
	s.logger.Info("Payment processed successfully", zap.Uint("paymentID", paymentID))
	return nil
}

func (s *CreditPaymentService) createInstallments(ctx context.Context, payment *domains.Payment) ([]domains.Installment, error) {
	s.logger.Info("Creating installments for payment", zap.Uint("paymentID", payment.ID))
	
	numberOfInstallments := s.calculateNumberOfInstallments(payment.Amount)
	installmentAmount := payment.Amount / numberOfInstallments
	
	var installments []domains.Installment
	for i := 1; i <= numberOfInstallments; i++ {
		installment := domains.Installment{
			PaymentID:         payment.ID,
			InstallmentNumber: i,
			DueDate:           time.Now().AddDate(0, i, 0).Format("2006-01-02"),
			Amount:            installmentAmount,
			Status:            "PENDING",
		}
		if i == numberOfInstallments {
			// Adjust the last installment to account for any rounding errors
			installment.Amount += payment.Amount - (installmentAmount * numberOfInstallments)
		}
		
		err := s.installmentRepo.Create(ctx, &installment)
		if err != nil {
			s.logger.Error("Failed to create installment", zap.Error(err))
			return nil, fmt.Errorf("failed to create installment: %w", err)
		}
		installments = append(installments, installment)
	}
	
	s.logger.Info("Installments created successfully", zap.Int("count", len(installments)))
	return installments, nil
}

func (s *CreditPaymentService) calculateNumberOfInstallments(amount int) int {
	// This is a simplified example. In a real-world scenario, this could be more complex
	// and might depend on various factors like credit score, payment amount, etc.
	if amount < 1000 {
		return 3
	} else if amount < 5000 {
		return 6
	} else {
		return 12
	}
}

func (s *CreditPaymentService) ProcessInstallment(ctx context.Context, installmentID uint) error {
	s.logger.Info("Processing installment", zap.Uint("installmentID", installmentID))
	
	installment, err := s.installmentRepo.GetByID(ctx, installmentID)
	if err != nil {
		s.logger.Error("Failed to get installment", zap.Error(err))
		return fmt.Errorf("failed to get installment: %w", err)
	}
	
	// Instead of processing the installment immediately, we'll just mark it as pending
	installment.Status = "PENDING"
	err = s.installmentRepo.Update(ctx, installment)
	if err != nil {
		s.logger.Error("Failed to update installment status", zap.Error(err))
		return fmt.Errorf("failed to update installment status: %w", err)
	}
	
	// Simulate a webhook
	go func() {
		time.Sleep(5 * time.Second) // Simulate a delay
		err := s.paymentGateway.SimulateInstallmentWebhook(context.Background(), installmentID)
		if err != nil {
			s.logger.Error("Failed to simulate installment webhook", zap.Error(err))
		}
	}()
	
	s.logger.Info("Installment marked as pending", zap.Uint("installmentID", installmentID))
	return nil
}

func (s *CreditPaymentService) HandleInstallmentWebhook(ctx context.Context, installmentID uint, status string) error {
	s.logger.Info("Handling installment webhook", zap.Uint("installmentID", installmentID), zap.String("status", status))
	
	installment, err := s.installmentRepo.GetByID(ctx, installmentID)
	if err != nil {
		s.logger.Error("Failed to get installment", zap.Error(err))
		return fmt.Errorf("failed to get installment: %w", err)
	}
	
	if status == "PAID" {
		installment.Status = "PAID"
	} else if status == "FAILED" {
		installment.Status = "FAILED"
	} else {
		return fmt.Errorf("unknown installment status: %s", status)
	}
	
	err = s.installmentRepo.Update(ctx, installment)
	if err != nil {
		s.logger.Error("Failed to update installment", zap.Error(err))
		return fmt.Errorf("failed to update installment: %w", err)
	}
	
	s.logger.Info("Installment webhook handled successfully", zap.Uint("installmentID", installmentID), zap.String("status", status))
	return nil
}

func (s *CreditPaymentService) GetPaymentDetails(ctx context.Context, paymentID uint) (*domains.Payment, error) {
    s.logger.Info("Getting payment details", zap.Uint("paymentID", paymentID))
    
    payment, err := s.paymentRepo.GetByID(ctx, paymentID)
    if err != nil {
        s.logger.Error("Failed to get payment", zap.Error(err))
        return nil, fmt.Errorf("failed to get payment: %w", err)
    }
    
    installmentPointers, err := s.installmentRepo.GetByPaymentID(ctx, paymentID)
    if err != nil {
        s.logger.Error("Failed to get installments for payment", zap.Error(err))
        return nil, fmt.Errorf("failed to get installments for payment: %w", err)
    }
    
    // Convert []*domains.Installment to []domains.Installment
    installments := make([]domains.Installment, len(installmentPointers))
    for i, installmentPtr := range installmentPointers {
        if installmentPtr != nil {
            installments[i] = *installmentPtr
        }
    }
    
    payment.Installments = installments
    return payment, nil
}
func (s *CreditPaymentService) getTotalPaymentsForCreditApp(ctx context.Context, creditAppID uint) (int, error) {
	payments, err := s.paymentRepo.GetByCreditApplicationID(ctx, creditAppID)
	if err != nil {
		return 0, err
	}
	
	total := 0
	for _, payment := range payments {
		if payment.Status == "SUCCESSFUL" {
			total += payment.Amount
		}
	}
	
	return total, nil
}