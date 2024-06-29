package credit

// CreditScoreRepository defines the interface for credit score data access
type CreditScoreRepository interface {
	Save(score *CreditScore) error
	GetByUserID(userID string) (*CreditScore, error)
	Update(score *CreditScore) error
}

// CreditAssessmentRepository defines the interface for credit assessment data access
type CreditAssessmentRepository interface {
	Save(assessment *CreditAssessment) error
	GetByUserID(userID string) (*CreditAssessment, error)
}

// CreditScoreCalculator defines the interface for calculating credit scores
type CreditScoreCalculator interface {
	CalculateScore(userID string, criteria CreditCriteria) (*CreditScore, error)
}

// CreditAssessmentService defines the interface for performing credit assessments
type CreditAssessmentService interface {
	AssessCredit(application *CreditApplication, criteria CreditCriteria) (*CreditAssessment, error)
}
