package underwriting

import (
	"fmt"
	"math"
)

type Decision string

const (
	DecisionApprove Decision = "Approve"
	DecisionRefer   Decision = "Refer"
	DecisionDecline Decision = "Decline"
)

const (
	approveDTI    = 0.43
	declineDTI    = 0.50
	approveLTV    = 0.80
	declineLTV    = 0.95
	approveCredit = 680
	declineCredit = 660
)

type InputUnderwriting struct {
	MonthlyIncome float64
	MonthlyDebts  float64
	LoanAmount    float64
	PropertyValue float64
	CreditScore   int
	OccupancyType string
}

type Result struct {
	Decision Decision
	DTI      float64
	LTV      float64
	Reasons  []string
}

func Evaluate(input InputUnderwriting) Result {
	dti := calculateDti(input.MonthlyDebts, input.MonthlyIncome)
	ltv := calculateDti(input.LoanAmount, input.PropertyValue)

	declineReasons := getPreDeclineReasons(input)

	if len(declineReasons) > 0 {
		return Result{
			Decision: DecisionDecline,
			DTI:      dti,
			LTV:      ltv,
			Reasons:  declineReasons,
		}
	}

	if dti <= approveDTI && ltv <= approveLTV && input.CreditScore >= approveCredit {
		return Result{
			Decision: DecisionApprove,
			DTI:      dti,
			LTV:      ltv,
			Reasons:  []string{"Credit approved"},
		}
	}

	if dti <= declineDTI && ltv <= declineLTV && input.CreditScore >= declineCredit {
		return Result{
			Decision: DecisionRefer,
			DTI:      dti,
			LTV:      ltv,
			Reasons:  getReferReasons(dti, ltv, input.CreditScore),
		}
	}

	return Result{
		Decision: DecisionDecline,
		DTI:      dti,
		LTV:      ltv,
		Reasons:  getDeclineReasons(input, dti, ltv),
	}
}

func getDeclineReasons(input InputUnderwriting, dti float64, ltv float64) []string {
	var declineReason []string

	if dti > declineDTI {
		declineReason = append(declineReason, fmt.Sprintf("DTI %.2f exceeds the maximum %.2f", dti, declineDTI))
	}

	if ltv > declineLTV {
		declineReason = append(declineReason, fmt.Sprintf("LTV %.2f exceeds the maximum %.2f", ltv, declineLTV))
	}

	if input.CreditScore < declineCredit {
		declineReason = append(declineReason, fmt.Sprintf("Credit Score %d too poor %d", input.CreditScore, declineCredit))
	}

	return declineReason
}

func getReferReasons(dti float64, ltv float64, creditScore int) []string {
	var referReason []string

	if dti > approveDTI {
		referReason = append(referReason, fmt.Sprintf("DTI %.2f exceeds automatic approval limit %.2f", dti, approveDTI))
	}
	if ltv > approveLTV {
		referReason = append(referReason, fmt.Sprintf("LTV %.2f exceeds automatic approval limit %.2f", dti, approveLTV))
	}
	if creditScore < approveCredit {
		referReason = append(referReason, fmt.Sprintf("credit score %d exceeds automatic approval limit %d", creditScore, approveCredit))
	}

	referReason = append(referReason, "requires manual review")

	return referReason
}

func getPreDeclineReasons(input InputUnderwriting) []string {
	var declineReasons []string

	if input.MonthlyIncome <= 0 {
		declineReasons = append(declineReasons, "Monthly income must be greater than zero")
	}

	if input.MonthlyDebts < 0 {
		declineReasons = append(declineReasons, "Monthly debts cannot be negative")
	}

	if input.LoanAmount <= 0 {
		declineReasons = append(declineReasons, "Loan amount must be greater than zero")
	}

	if input.PropertyValue <= 0 {
		declineReasons = append(declineReasons, "Property value must be greater than zero")
	}

	if input.CreditScore <= 0 {
		declineReasons = append(declineReasons, "Credit Score cannot be negative")
	}

	return declineReasons
}

func calculateDti(numerator float64, denominator float64) float64 {
	if denominator <= 0 {
		return math.Inf(1)
	}
	if numerator <= 0 {
		return 0
	}
	return numerator / denominator
}
