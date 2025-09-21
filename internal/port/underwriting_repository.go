package port

import (
	domain "AgilityFeat-Backend/internal/core/underwriting"
	"context"
)

type EvaluationRecord struct {
	UserID      string
	Application domain.InputUnderwriting
	Result      domain.Result
}

type UnderwritingRepository interface {
	Save(ctx context.Context, record EvaluationRecord) error
	FindByUser(ctx context.Context, userID string) ([]EvaluationRecord, error)
}
