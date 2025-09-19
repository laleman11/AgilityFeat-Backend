package underwriting

import (
	"context"
	"errors"
	"testing"
)

func TestServiceEvaluate(t *testing.T) {
	svc := NewService()

	result, err := svc.Evaluate(context.Background(), ApplicationInput{
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
	svc := NewService()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := svc.Evaluate(ctx, ApplicationInput{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}
