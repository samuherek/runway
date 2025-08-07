package engine

import (
	"fmt"
	"time"
)

// type ScenarioSettings struct {
// 	BaseInflationRate       float64
// 	InflationSpikeChance    float64
// 	InflationSpikeMin       float64
// 	InflationSpikeMax       float64
// 	UnexpectedExpenseChance float64
// 	UnexpectedExpenseMin    float64
// 	UnexpectedExpenseMax    float64
// 	Years                   int
// }

type SimpleRetirementInput struct {
	MonthlyExpense    float64
	YearsToRetirement int
	YearsInRetirement int
	Cash              float64
	Income            float64
}

func (s *SimpleRetirementInput) IntoSimulationState() SimulationState {
	today := time.Now()
	totalMonths := (s.YearsToRetirement + s.YearsInRetirement) * 12
	inflations := generateMonthlyInflations(totalMonths)

	bankAccount := Asset{
		ID:                 generateID(),
		Name:               "Bank account",
		Kind:               BankAccount,
		Amount:             s.Cash,
		StartMonthIndex:    0,
		EndMonthIndex:      nil,
		AnnualChangeChance: 0,
		AnnualChangeMin:    0,
		AnnualChangeMax:    0,
	}

	return SimulationState{
		MonthIndex: 0,
		StartYear:  today.Year(),
		StartMonth: int(today.Month()),
		Assets:     []Asset{bankAccount},
		Incomes: []Income{{
			ID:                 generateID(),
			Name:               "Work",
			Amount:             s.Income,
			StartMonthIndex:    0,
			EndMonthIndex:      nil,
			AnnualChangeChance: 0,
			AnnualChangeMin:    0,
			AnnualChangeMax:    0,
			AssetLinks: []DestinationLink{{
				TargetKind: EntityAssetKind,
				TargetID:   bankAccount.ID,
				Proportion: 1.0,
			}},
		}},
		Expenses: []Expense{{
			ID:              generateID(),
			Name:            "Life",
			Amount:          s.MonthlyExpense,
			StartMonthIndex: 0,
			EndMonthIndex:   nil,
			Frequency:       Monthly,
			InflationLinked: true,
			AssetLinks: []DestinationLink{{
				TargetKind: EntityAssetKind,
				TargetID:   bankAccount.ID,
				Proportion: 1.0,
			}},
		}},
		Liabilities: []Liability{},
		Events:      []Event{},
		Inflations:  inflations,
		History:     []MonthlySnapshot{},
	}
}

type RetirementQueryResult struct {
	RequiredSavingsAtRetirement          float64
	RequiredSavingsTodayWithoutInflation float64
	CashAtRetireTime                     float64
	CashToRetireDiff                     float64
}

func QueryRetirementPlan(history []MonthlySnapshot, input SimpleRetirementInput) RetirementQueryResult {
	retirementMonth := input.YearsToRetirement * 12
	retirementDurationMonths := input.YearsInRetirement * 12
	retirementEnd := retirementMonth + retirementDurationMonths

	fmt.Printf("RES: %v\n", history[retirementMonth].Expense)

	retirementPrice := 0.0
	for offset := retirementMonth; offset < retirementEnd; offset++ {
		if offset >= 0 && offset < len(history) {
			retirementPrice += history[offset].Expense
		}
	}

	accInflationAtRetireTime := 1.0
	for i := 0; i < retirementMonth; i++ {
		if i < len(history) {
			accInflationAtRetireTime *= 1 + history[i].Inflation
		}
	}

	retirementPriceToday := retirementPrice / accInflationAtRetireTime

	totalIncomeAtRetire := 0.0
	totalExpenseAtRetire := 0.0
	totalNetDiffAtRetire := 0.0
	for i := 0; i < retirementMonth; i++ {
		if i < len(history) {
			totalIncomeAtRetire += history[i].Income
			totalExpenseAtRetire += history[i].Expense
			totalNetDiffAtRetire += history[i].NetChange
		}
	}

	cashAtRetireTime := history[retirementMonth].NetWorth
	cashToRetireDiff := cashAtRetireTime - retirementPrice

	fmt.Printf("Total income at retire: %.2f\n", totalIncomeAtRetire)
	fmt.Printf("Total expense at retire: %.2f\n", totalExpenseAtRetire)
	fmt.Printf("Total retire price: %.2f\n", retirementPrice)
	fmt.Printf("Net worth at retire: %.2f\n", history[retirementMonth].NetWorth)
	fmt.Printf("Cash diff: %.2f\n", cashToRetireDiff)

	return RetirementQueryResult{
		RequiredSavingsAtRetirement:          retirementPrice,
		RequiredSavingsTodayWithoutInflation: retirementPriceToday, // since no growth, this is the same
		CashAtRetireTime:                     cashAtRetireTime,
		CashToRetireDiff:                     cashToRetireDiff,
	}
}

// func QueryRetirementPlan(history []MonthlySnapshot, input SimpleRetirementInput) RetirementQueryResult {
// 	retirementMonth := input.YearsToRetirement * 12
// 	endMonth := retirementMonth + input.YearsInRetirement*12
//
// 	// Step 1: Calculate total retirement need at retirement start
// 	// Present value of 20 years of expenses, inflation-adjusted
// 	monthlyRealExpense := input.MonthlyExpense * math.Pow(1+input.InflationRate, float64(input.YearsToRetirement))
// 	monthlyRate := input.ExpectedReturnRate / 12
// 	n := float64(input.YearsInRetirement * 12)
//
// 	requiredSavings := monthlyRealExpense * ((1 - math.Pow(1+monthlyRate, -n)) / monthlyRate)
//
// 	// Step 2: How much did the user accumulate by retirement month?
// 	savingsAchieved := sumAssetValueAtMonth(history, retirementMonth)
//
// 	// Step 3: Back-calculate what this amount would be today (discounted)
// 	requiredSavingsToday := requiredSavings / math.Pow(1+input.InflationRate, float64(input.YearsToRetirement))
//
// 	// Step 4: If short, calculate how much more needs to be saved monthly
// 	shortfall := requiredSavings - savingsAchieved
//
// 	additionalMonthly := 0.0
// 	if shortfall > 0 {
// 		additionalMonthly = shortfall / n // crude approximation
// 	}
//
// 	return RetirementQueryResult{
// 		RequiredSavingsAtRetirement:          requiredSavings,
// 		RequiredSavingsTodayWithoutInflation: requiredSavingsToday,
// 		SavingsAchieved:                      savingsAchieved,
// 		Shortfall:                            shortfall,
// 		AdditionalMonthlySavingsRequired:     additionalMonthly,
// 		CanReachGoal:                         shortfall <= 0,
// 	}
// }

// func sumAssetValueAtMonth(history []MonthlySnapshot, month int) float64 {
// 	if month >= len(history) {
// 		return 0
// 	}
// 	snapshot := history[month]
// 	sum := 0.0
// 	for _, asset := range snapshot.Assets {
// 		sum += asset.Value
// 	}
// 	return sum
// }

// func FindMinimumMonthlySavings(input RetirementInput) float64 {
// 	for add := 0.0; add < 5000; add += 50 {
// 		inputWithExtra := input
// 		inputWithExtra.MonthlyExpenses -= add
//
// 		state := BuildInitialSimulation(inputWithExtra)
// 		RunSimulationLoop(state, (input.YearsToRetirement+input.YearsInRetirement)*12)
// 		result := QueryRetirementPlan(state.History, input)
//
// 		if result.CanReachGoal {
// 			return add
// 		}
// 	}
// 	return -1 // not reachable
// }
