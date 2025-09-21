package underwriting

import (
	"AgilityFeat-Backend/internal/infra/memory/underwriting"
	"AgilityFeat-Backend/internal/port"
	"context"
	"errors"
	"testing"
)

func TestServiceEvaluate(t *testing.T) {
	repository := underwriting.NewRepository()
	svc := NewService(repository)

	result, err := svc.Evaluate(context.Background(), port.UnderwritingRequest{
		MonthlyIncome: 8000,
		MonthlyDebts:  2000,
		LoanAmount:    250000,
		PropertyValue: 350000,
		CreditScore:   760,
		OccupancyType: "primary",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Decision != "Approve" {
		t.Fatalf("expected approve, got %s", result.Decision)
	}
}

func TestServiceEvaluateContextCanceled(t *testing.T) {
	repository := underwriting.NewRepository()
	svc := NewService(repository)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := svc.Evaluate(ctx, port.UnderwritingRequest{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}
