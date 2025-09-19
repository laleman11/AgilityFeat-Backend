package underwriting

import (
	underwritingCore "AgilityFeat-Backend/internal/core/underwriting"
	"context"
)

type ApplicationInput struct {
	MonthlyIncome float64
	MonthlyDebts  float64
	LoanAmount    float64
	PropertyValue float64
	CreditScore   int
	OccupancyType string
}

type Result struct {
	Decision string
	DTI      float64
	LTV      float64
	Reasons  []string
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Evaluate(ctx context.Context, input ApplicationInput) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	result := underwritingCore.Evaluate(underwritingCore.InputUnderwriting{
		MonthlyIncome: input.MonthlyIncome,
		MonthlyDebts:  input.MonthlyDebts,
		LoanAmount:    input.LoanAmount,
		PropertyValue: input.PropertyValue,
		CreditScore:   input.CreditScore,
		OccupancyType: input.OccupancyType,
	})

	return Result{
		Decision: string(result.Decision),
		DTI:      result.DTI,
		LTV:      result.LTV,
		Reasons:  result.Reasons,
	}, nil
}
