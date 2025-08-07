package engine

import (
	"math"
	"time"
)

type RetirementInput struct {
	MonthlyExpenseToday  float64
	MonthlyIncome        float64
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
	RequiredMonthlySaving         float64   // how much to save monthly to meet the goal
	SavingProgression             []float64 // growing monthly savings over time
	ReachedTargetInMonths         int       // how many months it would take with current income-expense saving
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

	// Calculate required monthly savings to reach target in time (adjusted for inflation)
	monthsUntilRetirement := input.YearsUntilRetirement * 12
	missingAmount := totalPresentValue - input.CurrentSavings
	requiredMonthlySaving := 0.0
	if missingAmount > 0 && monthsUntilRetirement > 0 {
		monthlyRate := math.Pow(1+r, 1.0/12.0) - 1
		requiredMonthlySaving = missingAmount * monthlyRate / (math.Pow(1+monthlyRate, float64(monthsUntilRetirement)) - 1)
	}

	// Build the growing savings plan
	savingProgression := make([]float64, monthsUntilRetirement)
	for i := 0; i < monthsUntilRetirement; i++ {
		// savingProgression[i] = requiredMonthlySaving * math.Pow(1+monthlyInflation, float64(i))
		savingProgression[i] = requiredMonthlySaving * math.Pow(1+r, float64(i)/12.0)
	}

	// Calculate how long it would take to reach the goal using (income - expenses) as monthly savings
	surplus := input.MonthlyIncome - input.MonthlyExpenseToday
	reachedMonth := -1
	if surplus > 0 {
		presentValueSum := input.CurrentSavings
		for i := 0; i < monthsUntilRetirement; i++ {
			// discount each month’s surplus individually
			discounted := surplus / math.Pow(1+r, float64(i)/12.0)
			presentValueSum += discounted
			if presentValueSum >= totalPresentValue {
				reachedMonth = i + 1
				break
			}
		}
	}

	return RetirementProjection{
		MonthlyValues:                 expenses,
		MonthlyExpenseAtStart:         monthlyExpenseAtStart,
		TotalRequiredFutureFund:       total,
		TotalRequiredPresentValueFund: totalPresentValue,
		CurrentCoveragePercentage:     coverage,
		RequiredMonthlySaving:         requiredMonthlySaving,
		SavingProgression:             savingProgression,
		ReachedTargetInMonths:         reachedMonth,
	}

}

type FundPoint struct {
	Date  time.Time
	Value float64
	Sum   float64
}

type RetirementProjectionData struct {
	CurrentMonthlyExpenses float64
	CurrentMonthlyIncome   float64
	YearsUntilRetirement   int
	YearsInRetirement      int
	CurrentSavings         float64
	InflationRate          float64
	// -
	MonthlyExpenseAtStart      float64
	RequiredFutureFund         float64
	RequiredPresentFund        float64
	CurrentFullfilmentPercent  float64
	FulfilmentMonths           int
	FutureFundWithdrawals      []FundPoint
	CurrentSavingsProgressions []FundPoint
}

func float2Clamp(value float64) float64 {
	return math.Round(value*100) / 100
}

func RetirementProjectionResult(input RetirementInput, projection RetirementProjection) RetirementProjectionData {
	var data RetirementProjectionData

	data.CurrentMonthlyExpenses = float2Clamp(input.MonthlyExpenseToday)
	data.CurrentMonthlyIncome = float2Clamp(input.MonthlyIncome)
	data.YearsUntilRetirement = input.YearsUntilRetirement
	data.YearsInRetirement = input.WithdrawalYears
	data.CurrentSavings = float2Clamp(input.CurrentSavings)
	data.InflationRate = float2Clamp(input.InflationRate)

	data.MonthlyExpenseAtStart = float2Clamp(projection.MonthlyExpenseAtStart)
	data.RequiredFutureFund = float2Clamp(projection.TotalRequiredFutureFund)
	data.RequiredPresentFund = float2Clamp(projection.TotalRequiredPresentValueFund)
	data.CurrentFullfilmentPercent = float2Clamp(projection.CurrentCoveragePercentage)
	data.FulfilmentMonths = projection.ReachedTargetInMonths

	rest := projection.TotalRequiredFutureFund
	for i, v := range projection.MonthlyValues {
		rest -= v
		year := data.YearsUntilRetirement + (i / 12)
		month := i % 12
		data.FutureFundWithdrawals = append(data.FutureFundWithdrawals, FundPoint{
			Date:  time.Now().AddDate(year, month, 0),
			Value: float2Clamp(v),
			Sum:   float2Clamp(rest),
		})
	}

	rest2 := input.CurrentSavings
	for i, v := range projection.SavingProgression {
		rest2 += v
		year := i / 12
		month := i % 12
		data.CurrentSavingsProgressions = append(data.CurrentSavingsProgressions, FundPoint{
			Date:  time.Now().AddDate(year, month, 0),
			Value: float2Clamp(v),
			Sum:   float2Clamp(rest2),
		})
	}

	return data
}
