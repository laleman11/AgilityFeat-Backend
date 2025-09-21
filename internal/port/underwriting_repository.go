package port

import (
	domain "AgilityFeat-Backend/internal/core/underwriting"
	"context"
	"time"
)

type EvaluationRecord struct {
	UserID      string
	Application domain.InputUnderwriting
	Result      domain.Result
	CreatedAt   time.Time
}

type UnderwritingRepository interface {
	Save(ctx context.Context, record EvaluationRecord) error
	FindByUser(ctx context.Context, userID string) ([]EvaluationRecord, error)
}
