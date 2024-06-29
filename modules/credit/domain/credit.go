package credit

import (
    "time"
)

// CreditScore represents a user's credit score
type CreditScore struct {
    UserID      string
    Score       int
    ScoreDate   time.Time
    RiskCategory string
    ScoreFactors []string
}

// CreditAssessment represents the result of a credit assessment
type CreditAssessment struct {
    UserID           string
    AssessmentDate   time.Time
    CriteriaMet      map[string]bool
    AssessmentResult string
    RequestedAmount  float64
    ApprovedAmount   float64
    ExpirationDate   time.Time
}

// CreditCriteria represents the criteria used for credit assessment
type CreditCriteria struct {
    IncomeLevel       float64
    DebtLevel         float64
    PaymentHistory    string
    EmploymentStatus  string
    CreditUtilization float64
    CreditHistoryLength int // in months
}

// CreditApplication represents a user's application for credit
type CreditApplication struct {
    ApplicationID string
    UserID        string
    RequestedAmount float64
    Purpose       string
    ApplicationDate time.Time
    Status        string // "Pending", "Approved", "Rejected"
}

// CreditLimit represents a user's credit limit
type CreditLimit struct {
    UserID        string
    CurrentLimit  float64
    AvailableCredit float64
    LastUpdateDate time.Time
}