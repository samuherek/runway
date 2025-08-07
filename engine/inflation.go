package engine

import "math"

func generateMonthlyInflations(months int) []Inflation {
	inflations := []Inflation{}
	monthlyRate := math.Pow(1+0.03, 1.0/12.0) - 1

	for i := 0; i < months; i++ {
		inflations = append(inflations, Inflation{
			Month: i,
			Value: monthlyRate,
		})
	}

	return inflations
}

// func generateInflationAndShocks(settings ScenarioSettings) ([]float64, map[int]float64, map[int]float64) {
// 	rand.Seed(time.Now().UnixNano())
// 	years := settings.Years
// 	inflation := make([]float64, years)
// 	incomeDrops := make(map[int]float64)
// 	extraExpenses := make(map[int]float64)
//
// 	for i := 0; i < years; i++ {
// 		inf := settings.BaseInflationRate
// 		if rand.Float64() < settings.InflationSpikeChance {
// 			inf += settings.InflationSpikeMin + rand.Float64()*(settings.InflationSpikeMax-settings.InflationSpikeMin)
// 		}
// 		inflation[i] = inf
//
// 		if rand.Float64() < settings.DropInIncomeChance {
// 			incomeDrops[i] = settings.DropPercentage
// 		}
// 		if rand.Float64() < settings.UnexpectedExpenseChance {
// 			extra := settings.UnexpectedExpenseMin + rand.Float64()*(settings.UnexpectedExpenseMax-settings.UnexpectedExpenseMin)
// 			extraExpenses[i] = extra
// 		}
// 	}
// 	return inflation, incomeDrops, extraExpenses
// }
