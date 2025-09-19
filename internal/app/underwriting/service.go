package underwriting

import "context"

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

	//TODO: fix
	return Result{
		Decision: "Approve",
		DTI:      2.2,
		LTV:      2.2,
		Reasons:  make([]string, 0),
	}, nil
}
