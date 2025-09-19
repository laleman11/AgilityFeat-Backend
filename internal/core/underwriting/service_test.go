package underwriting

import "testing"

func TestEvaluateApprove(t *testing.T) {
	result := Evaluate(InputUnderwriting{
		MonthlyIncome: 10000,
		MonthlyDebts:  4000,
		LoanAmount:    240000,
		PropertyValue: 300000,
		CreditScore:   700,
	})

	if result.Decision != DecisionApprove {
		t.Fatalf("expected approve, got %s", result.Decision)
	}

	if len(result.Reasons) != 1 || result.Reasons[0] != "Credit approved" {
		t.Fatalf("unexpected reasons: %v", result.Reasons)
	}
}

func TestEvaluateRefer(t *testing.T) {
	result := Evaluate(InputUnderwriting{
		MonthlyIncome: 8000,
		MonthlyDebts:  3600,
		LoanAmount:    260000,
		PropertyValue: 300000,
		CreditScore:   670,
	})

	if result.Decision != DecisionRefer {
		t.Fatalf("expected refer, got %s", result.Decision)
	}

	if len(result.Reasons) == 0 {
		t.Fatal("expected reasons for refer decision")
	}
}

func TestEvaluateDeclineThresholds(t *testing.T) {
	result := Evaluate(InputUnderwriting{
		MonthlyIncome: 7000,
		MonthlyDebts:  3800,
		LoanAmount:    320000,
		PropertyValue: 320000,
		CreditScore:   640,
	})

	if result.Decision != DecisionDecline {
		t.Fatalf("expected decline, got %s", result.Decision)
	}

	if len(result.Reasons) == 0 {
		t.Fatal("expected decline reasons")
	}
}

func TestEvaluateDeclineForInvalidInputs(t *testing.T) {
	result := Evaluate(InputUnderwriting{})

	if result.Decision != DecisionDecline {
		t.Fatalf("expected decline, got %s", result.Decision)
	}

	if len(result.Reasons) < 3 {
		t.Fatalf("expected multiple reasons for invalid input, got %v", result.Reasons)
	}
}
