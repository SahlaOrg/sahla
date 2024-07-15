package utils

import (
	"math/rand"
	"time"
)

// CreditScoreRange defines the range for credit scores
type CreditScoreRange struct {
	Min int
	Max int
}

var (
	PoorScore      = CreditScoreRange{300, 579}
	FairScore      = CreditScoreRange{580, 669}
	GoodScore      = CreditScoreRange{670, 739}
	VeryGoodScore  = CreditScoreRange{740, 799}
	ExcellentScore = CreditScoreRange{800, 850}
)

// GenerateRandomCreditScore generates a random credit score between 300 and 850
func GenerateRandomCreditScore() int {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Define the overall range for credit scores
	minScore := 300
	maxScore := 850

	// Generate a random score within the range
	score := rand.Intn(maxScore-minScore+1) + minScore

	return score
}

// GenerateWeightedRandomCreditScore generates a random credit score with a weighted distribution
// This function is more likely to generate scores in the "Fair" and "Good" ranges
func GenerateWeightedRandomCreditScore() int {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Define weights for each credit score range
	weights := []int{
		10,  // Poor
		25,  // Fair
		35,  // Good
		20,  // Very Good
		10,  // Excellent
	}

	// Calculate the total weight
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	// Generate a random number within the total weight
	randomWeight := rand.Intn(totalWeight)

	// Determine which range the random number falls into
	var selectedRange CreditScoreRange
	for i, weight := range weights {
		if randomWeight < weight {
			switch i {
			case 0:
				selectedRange = PoorScore
			case 1:
				selectedRange = FairScore
			case 2:
				selectedRange = GoodScore
			case 3:
				selectedRange = VeryGoodScore
			case 4:
				selectedRange = ExcellentScore
			}
			break
		}
		randomWeight -= weight
	}

	// Generate a random score within the selected range
	score := rand.Intn(selectedRange.Max-selectedRange.Min+1) + selectedRange.Min

	return score
}