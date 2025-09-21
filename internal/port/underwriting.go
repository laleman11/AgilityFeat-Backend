package port

import (
	"context"
)

type UnderwritingRequest struct {
	UserID        string
	MonthlyIncome float64
	MonthlyDebts  float64
	LoanAmount    float64
	PropertyValue float64
	CreditScore   int
	OccupancyType string
}

type UnderwritingResponse struct {
	Decision string
	DTI      float64
	LTV      float64
	Reasons  []string
}

type UnderwritingHistoryEntry struct {
	UserID   string
	Request  UnderwritingRequest
	Response UnderwritingResponse
}

type UnderWritingService interface {
	Evaluate(ctx context.Context, input UnderwritingRequest) (UnderwritingResponse, error)
	History(ctx context.Context, userID string) ([]UnderwritingHistoryEntry, error)
}
