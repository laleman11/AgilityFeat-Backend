package port

import (
	underwriting "AgilityFeat-Backend/internal/app/underwriting"
	"context"
)

type UnderWritingService interface {
	Evaluate(ctx context.Context, input underwriting.ApplicationInput) (underwriting.Result, error)
}
