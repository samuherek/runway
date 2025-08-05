package engine

import (
	"math"
)

type RetirementInput struct {
	MonthlyExpenseToday  float64
	YearsUntilRetirement int
	InflationRate        float64
	WithdrawalYears      int
	CurrentSavings       float64
}

type RetirementProjection struct {
	MonthlyValues                 []float64 // inflated expenses per month for 20 years
	MonthlyExpenseAtStart         float64   // what monthly expense will be at retirement
	TotalRequiredFutureFund       float64   // how much you need at retirement to cover 20 years
	TotalRequiredPresentValueFund float64   // how much that future total is worth in today's money
	CurrentCoveragePercentage     float64   // how much of the required present value is already saved
}

func ProjectRetirement(input RetirementInput) RetirementProjection {
	n := input.WithdrawalYears * 12
	years := float64(input.YearsUntilRetirement)
	r := input.InflationRate
	monthlyInflation := math.Pow(1+r, 1.0/12.0)

	// Adjust today's expense to retirement year
	monthlyExpenseAtStart := input.MonthlyExpenseToday * math.Pow(1+r, years)

	// Build inflation-adjusted expenses for each month
	expenses := make([]float64, n)
	total := 0.0
	for i := 0; i < n; i++ {
		adjusted := monthlyExpenseAtStart * math.Pow(monthlyInflation, float64(i))
		expenses[i] = adjusted
		total += adjusted
	}

	// Discount future total back to today's value
	discountFactor := math.Pow(1+r, years)
	totalPresentValue := total / discountFactor

	// Calculate percentage coverage by current savings
	coverage := 0.0
	if totalPresentValue > 0 {
		coverage = (input.CurrentSavings / totalPresentValue) * 100.0
	}

	return RetirementProjection{
		MonthlyValues:                 expenses,
		MonthlyExpenseAtStart:         monthlyExpenseAtStart,
		TotalRequiredFutureFund:       total,
		TotalRequiredPresentValueFund: totalPresentValue,
		CurrentCoveragePercentage:     coverage,
	}
}
