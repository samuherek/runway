package engine

import (
	"math"
	"math/rand"
	"time"
)

type OneTimeExpense struct {
	MonthOffset int
	Amount      float64
	Description string
}

type UnexpectedExpense struct {
	Probability float64
	MaxCost     float64
	Frequency   int
}

type InflationForecast struct {
	Year int
	Rate float64
}

type Input struct {
	InitialBalance     float64
	MonthlyIncome      float64
	MonthlyExpenses    float64
	ExpensesConfidence float64 // e.g., 0.9 for 90%
	PlannedExpenses    []OneTimeExpense
	UnexpectedExpense  UnexpectedExpense
	InflationRates     []InflationForecast
	MaxMonths          int
	Simulations        int
}

type ProjectionResult struct {
	Month    int
	Balances []float64
}

func getInflationRate(year int, forecast []InflationForecast) float64 {
	for _, f := range forecast {
		if f.Year == year {
			return f.Rate
		}
	}
	return 0.03
}

func SimulateSimpleProjection(input Input) []ProjectionResult {
	rand.Seed(time.Now().UnixNano())
	results := make([]ProjectionResult, input.MaxMonths)

	for sim := 0; sim < input.Simulations; sim++ {
		balance := input.InitialBalance
		monthly := make([]float64, input.MaxMonths)

		// Decide if unexpected expense happens
		unexpectedHappens := rand.Float64() < input.UnexpectedExpense.Probability
		unexpectedMonth := 0
		if unexpectedHappens {
			unexpectedMonth = rand.Intn(input.MaxMonths)
		}

		for month := 0; month < input.MaxMonths; month++ {
			year := 2025 + month/12
			inflation := getInflationRate(year, input.InflationRates)
			factor := math.Pow(1+inflation, float64(month)/12.0)

			income := input.MonthlyIncome / factor
			expenses := input.MonthlyExpenses / input.ExpensesConfidence * factor

			for _, e := range input.PlannedExpenses {
				if e.MonthOffset == month {
					expenses += e.Amount
				}
			}

			if unexpectedHappens && month == unexpectedMonth {
				expenses += input.UnexpectedExpense.MaxCost
			}

			balance += income - expenses
			monthly[month] = balance

			if balance <= 0 {
				break
			}
		}

		for i, v := range monthly {
			results[i].Month = i
			results[i].Balances = append(results[i].Balances, v)
		}
	}

	return results
}

type ProjectionMonth struct {
	Date string
	Min  float64
	Mid  float64
	Max  float64
}

func ExtractMinMax(input []ProjectionResult) []ProjectionMonth {
	now := time.Now()
	results := []ProjectionMonth{}
	for _, r := range input {
		if len(r.Balances) == 0 {
			continue
		}
		min := r.Balances[0]
		max := r.Balances[0]
		sum := 0.0
		for _, b := range r.Balances {
			if b < min {
				min = b
			}
			if b > max {
				max = b
			}
			sum += b
		}
		avg := sum / float64(len(r.Balances))

		if min == 0 && avg == 0 && max == 0 {
			continue
		}

		results = append(results, ProjectionMonth{
			Date: now.AddDate(0, r.Month, 0).Format("2006-01-02"),
			Min:  math.Max(0, math.Round(min*100)/100),
			Mid:  math.Max(0, math.Round(avg*100)/100),
			Max:  math.Max(0, math.Round(max*100)/100),
		})
	}

	return results
}
