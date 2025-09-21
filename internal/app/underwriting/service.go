package underwriting

import (
	underwritingCore "AgilityFeat-Backend/internal/core/underwriting"
	"AgilityFeat-Backend/internal/port"
	"context"
	"time"
)

type Service struct {
	repository port.UnderwritingRepository
}

func NewService(repository port.UnderwritingRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Evaluate(ctx context.Context, input port.UnderwritingRequest) (port.UnderwritingResponse, error) {
	if err := ctx.Err(); err != nil {
		return port.UnderwritingResponse{}, err
	}
	application := underwritingCore.InputUnderwriting{
		MonthlyIncome: input.MonthlyIncome,
		MonthlyDebts:  input.MonthlyDebts,
		LoanAmount:    input.LoanAmount,
		PropertyValue: input.PropertyValue,
		CreditScore:   input.CreditScore,
		OccupancyType: input.OccupancyType,
	}
	result := underwritingCore.Evaluate(application)
	now := time.Now().UTC()

	s.repository.Save(ctx, port.EvaluationRecord{
		UserID:      input.UserID,
		Result:      result,
		Application: application,
		CreatedAt:   now,
	})

	return port.UnderwritingResponse{
		Decision: string(result.Decision),
		DTI:      result.DTI,
		LTV:      result.LTV,
		Reasons:  result.Reasons,
	}, nil
}

func (s *Service) History(ctx context.Context, userID string) ([]port.UnderwritingHistoryEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if s.repository == nil {
		return []port.UnderwritingHistoryEntry{}, nil
	}

	records, err := s.repository.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	history := make([]port.UnderwritingHistoryEntry, len(records))
	for i, rec := range records {
		history[i] = port.UnderwritingHistoryEntry{
			UserID: rec.UserID,
			Request: port.UnderwritingRequest{
				UserID:        rec.UserID,
				MonthlyIncome: rec.Application.MonthlyIncome,
				MonthlyDebts:  rec.Application.MonthlyDebts,
				LoanAmount:    rec.Application.LoanAmount,
				PropertyValue: rec.Application.PropertyValue,
				CreditScore:   rec.Application.CreditScore,
				OccupancyType: rec.Application.OccupancyType,
			},
			Response: port.UnderwritingResponse{
				Decision: string(rec.Result.Decision),
				DTI:      rec.Result.DTI,
				LTV:      rec.Result.LTV,
				Reasons:  append([]string(nil), rec.Result.Reasons...),
			},
			CreatedAt: rec.CreatedAt.Format(time.RFC3339Nano),
		}
	}

	return history, nil
}
