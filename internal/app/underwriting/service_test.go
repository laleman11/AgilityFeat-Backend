package underwriting

import (
	"AgilityFeat-Backend/internal/port"
	"context"
	"errors"
	"testing"
)

type stubRepository struct {
	records []port.EvaluationRecord
	err     error
}

func (s *stubRepository) Save(ctx context.Context, record port.EvaluationRecord) error {
	if s.err != nil {
		return s.err
	}
	s.records = append(s.records, record)
	return nil
}

func (s *stubRepository) FindByUser(ctx context.Context, userID string) ([]port.EvaluationRecord, error) {
	if s.err != nil {
		return nil, s.err
	}

	out := make([]port.EvaluationRecord, len(s.records))
	copy(out, s.records)
	return out, nil
}

func TestServiceEvaluate(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository)

	result, err := service.Evaluate(context.Background(), port.UnderwritingRequest{
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
	repository := &stubRepository{}
	service := NewService(repository)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := service.Evaluate(ctx, port.UnderwritingRequest{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}

func TestServiceHistory(t *testing.T) {
	repo := &stubRepository{}
	service := NewService(repo)

	_, err := service.Evaluate(context.Background(), port.UnderwritingRequest{
		UserID:        "user-1",
		MonthlyIncome: 9000,
		MonthlyDebts:  3000,
		LoanAmount:    200000,
		PropertyValue: 300000,
		CreditScore:   700,
		OccupancyType: "primary",
	})
	if err != nil {
		t.Fatalf("unexpected error evaluating: %v", err)
	}

	history, err := service.History(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error retrieving history: %v", err)
	}

	if len(history) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(history))
	}

	if history[0].Request.UserID != "user-1" {
		t.Fatalf("unexpected user id: %s", history[0].Request.UserID)
	}
}
