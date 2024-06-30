package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	domain "github.com/mohamed2394/sahla/modules/credit"
	repositories "github.com/mohamed2394/sahla/modules/credit/repository"
)

type CreditService struct {
	repo   *repositories.GormCreditRepository
	apiURL string
}

func NewCreditService(repo *repositories.GormCreditRepository, apiURL string) *CreditService {
	return &CreditService{
		repo:   repo,
		apiURL: apiURL,
	}
}

func (s *CreditService) SaveCreditScore(score *domain.CreditScore) error {
	return s.repo.SaveCreditScore(score)
}

func (s *CreditService) GetCreditScoreByUserID(userID string) (*domain.CreditScore, error) {
	return s.repo.GetCreditScoreByUserID(userID)
}

func (s *CreditService) SaveCreditFeatures(features *domain.CreditFeatures) (*domain.CreditScore, error) {
	// Save features to the database
	if err := s.repo.SaveCreditFeatures(features); err != nil {
		return nil, err
	}

	// Communicate with Flask API
	flaskResponse, err := sendCreditFeaturesToFlaskAPI(*features, s.apiURL)
	if err != nil {
		return nil, err
	}

	// Extract score and risk category from the response
	score, ok := flaskResponse["credit_score"].(float64)
	if !ok {
		return nil, errors.New("invalid credit score in Flask API response")
	}
	riskCategory, ok := flaskResponse["risk_category"].(string)
	if !ok {
		return nil, errors.New("invalid risk category in Flask API response")
	}

	creditScore := &domain.CreditScore{
		UserID:       features.UserID,
		Score:        int(score),
		RiskCategory: riskCategory,
		ScoreDate:    time.Now(),
		ScoreFactors: []string{},
	}

	// Save credit score to the database
	if err := s.repo.SaveCreditScore(creditScore); err != nil {
		return nil, err
	}

	return creditScore, nil
}

func (s *CreditService) GetCreditFeaturesByUserID(userID string) (*domain.CreditFeatures, error) {
	return s.repo.GetCreditFeaturesByUserID(userID)
}
func sendCreditFeaturesToFlaskAPI(features domain.CreditFeatures, apiURL string) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(features)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response: %v", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
