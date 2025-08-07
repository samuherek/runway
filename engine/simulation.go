package engine

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type AssetKind int
type Percent float64

const (
	Cash AssetKind = iota
	BankAccount
	RealEstate
	Depreciating
	Investment
)

type Asset struct {
	ID                 ID
	Name               string    // e.g., "Emergency Fund", "Toyota Camry"
	Kind               AssetKind // from the enum above
	Amount             float64   // The actual current value at this point
	StartMonthIndex    int       // when it's acquired (if immediate start month is 0)
	EndMonthIndex      *int      // optional: if asset is sold or disposed
	AnnualChangeChance Percent   // 0 if not changable. house, stock... will have 1
	AnnualChangeMin    Percent   // Min range of the anual change in percent
	AnnualChangeMax    Percent   // Max range of the anual change in percent
}

func (a *Asset) GetID() ID         { return a.ID }
func (a *Asset) GetName() string   { return a.Name }
func (a *Asset) GetValue() float64 { return a.Amount }

type EntityKind string

const (
	EntityAssetKind     EntityKind = "asset"
	EntityIncomeKind    EntityKind = "income"
	EntityExpenseKind   EntityKind = "expense"
	EntityLiabilityKind EntityKind = "liability"
)

type DestinationLink struct {
	TargetKind EntityKind // The type of the entity in the state
	TargetID   ID         // We assume it is always an asset at this point
	Proportion Percent
}

type Income struct {
	ID                 ID
	Name               string            // e.g., "Salary", "Dividends"
	Amount             float64           // The current month value of the income
	StartMonthIndex    int               // If 0 starts immediatelly. Less likely in the future
	EndMonthIndex      *int              // If exists it means it ends at some point
	AnnualChangeChance Percent           // if 0 we assume no change
	AnnualChangeMin    Percent           // if some, we adjust. Like salary bump each year
	AnnualChangeMax    Percent           // Same as above
	AssetLinks         []DestinationLink // We need to put the money into some destinations
}

func (a *Income) GetID() ID         { return a.ID }
func (a *Income) GetName() string   { return a.Name }
func (a *Income) GetValue() float64 { return a.Amount }

type FrequencyCycle int

const (
	Once FrequencyCycle = iota
	Monthly
	Quarterly
	BiYearly
	Yearly
)

type Expense struct {
	ID              ID
	Name            string            // "Food"
	Amount          float64           //
	StartMonthIndex int               // If 0 it starts immediatelly
	EndMonthIndex   *int              // if nil, it never ends
	Frequency       FrequencyCycle    // Defines how often and if it repeats
	InflationLinked bool              // food, or rent then linked. One time payment is not
	AssetLinks      []DestinationLink // We need to put the money into some destinations
}

func (a *Expense) GetID() ID         { return a.ID }
func (a *Expense) GetName() string   { return a.Name }
func (a *Expense) GetValue() float64 { return a.Amount }

type Liability struct {
	ID              ID
	Name            string            //
	Amount          float64           // The amount to pay
	StartMonthIndex int               // When the liability starts
	EndMonthIndex   *int              // If it ever ends then it has the end month
	Frequency       FrequencyCycle    // How often is the "Amount" payed
	AssetRef        *ID               // If it is linked to asset (for effects)
	AssetLinks      []DestinationLink // We need to put the money into some destinations
}

func (a *Liability) GetID() ID         { return a.ID }
func (a *Liability) GetName() string   { return a.Name }
func (a *Liability) GetValue() float64 { return a.Amount }

type FinanceView struct {
	ID    ID
	Name  string
	Value float64
}

type FinanceViewable interface {
	GetID() ID
	GetName() string
	GetValue() float64
}

type MonthlySnapshot struct {
	MonthIndex  int
	Month       int
	Year        int
	Income      float64
	Expense     float64
	NetChange   float64
	NetWorth    float64
	Assets      []FinanceView
	Incomes     []FinanceView
	Expenses    []FinanceView
	Liabilities []FinanceView
	Inflation   float64
}

// TODO:: Rework this to accept parameters which would then have custom functions somewhere else to use those configs.
type Event struct {
	ID      ID
	Name    string                            // e.g., "Selling house", "Selling car"...
	Trigger func(state *SimulationState) bool // Dynamic filter based on current state. Either "current monht == 1" or "house value at $xxx"
	Effect  func(state *SimulationState)      // The actual function to mutate the simulation state
}

type Inflation struct {
	Month int
	Value float64
}

type SimulationState struct {
	MonthIndex  int
	StartYear   int
	StartMonth  int
	Assets      []Asset
	Incomes     []Income
	Expenses    []Expense
	Liabilities []Liability
	Events      []Event
	Inflations  []Inflation
	History     []MonthlySnapshot
}

func (s *SimulationState) currMonth() int {
	return (s.MonthIndex + s.StartMonth) % 12
}

func (s *SimulationState) currYear() int {
	return s.StartYear + ((s.MonthIndex + s.StartMonth) / 12)
}

func (s *SimulationState) sumIncomes() float64 {
	total := 0.0
	for i := range s.Incomes {
		total += s.Incomes[i].Amount
	}

	return total
}

func (s *SimulationState) sumExpenses() float64 {
	total := 0.0
	for i := range s.Expenses {
		total += s.Expenses[i].Amount
	}

	return total
}

func (s *SimulationState) sumAssets() float64 {
	total := 0.0
	for i := range s.Assets {
		total += s.Assets[i].Amount
	}

	return total
}

func (s *SimulationState) assetFinanceViews() []FinanceView {
	list := []FinanceView{}
	for _, el := range s.Assets {
		list = append(list, intoFinanceView(&el))
	}

	return list
}

func (s *SimulationState) incomeFinanceViews() []FinanceView {
	list := []FinanceView{}
	for _, el := range s.Incomes {
		list = append(list, intoFinanceView(&el))
	}

	return list
}

func (s *SimulationState) expenseFinanceViews() []FinanceView {
	list := []FinanceView{}
	for _, el := range s.Expenses {
		list = append(list, intoFinanceView(&el))
	}

	return list
}

func (s *SimulationState) liabilityFinanceViews() []FinanceView {
	list := []FinanceView{}
	for _, el := range s.Liabilities {
		list = append(list, intoFinanceView(&el))
	}

	return list
}

func intoFinanceView[T FinanceViewable](e T) FinanceView {
	return FinanceView{
		ID:    e.GetID(),
		Name:  e.GetName(),
		Value: e.GetValue(),
	}
}

// TODO: There could be case when we forget to check the "END".
// Maybe we need to pass more to the function or do some other check.
func isActiveRange(monthIndex, checkMonthIndex int, checkEndMonthIndex *int) bool {
	started := checkMonthIndex <= monthIndex
	ended := checkEndMonthIndex != nil && *checkEndMonthIndex >= monthIndex

	return started && !ended
}

func isExpired(monthIndex int, checkEndMonthIndex *int) bool {
	return checkEndMonthIndex != nil && *checkEndMonthIndex < monthIndex
}

func (s *SimulationState) snapshot() MonthlySnapshot {
	income := s.sumIncomes()
	expense := s.sumExpenses()
	netWorth := s.sumAssets()
	inflation := s.Inflations[0].Value

	snap := MonthlySnapshot{
		MonthIndex:  s.MonthIndex,
		Month:       s.currMonth(),
		Year:        s.currYear(),
		Income:      income,
		Expense:     expense,
		NetChange:   income - expense,
		NetWorth:    netWorth,
		Assets:      s.assetFinanceViews(),
		Incomes:     s.incomeFinanceViews(),
		Expenses:    s.expenseFinanceViews(),
		Liabilities: s.liabilityFinanceViews(),
		Inflation:   inflation,
	}

	return snap
}

func applyAssetChange(asset *Asset, state *SimulationState) {
	// TODO: There are no really any changes here right? It's usually static asset.
	// But things like real estate or stock can have some predefined possible change.
	// Until then, we don't do anything here.
}

func applyIncomeChange(income *Income, state *SimulationState) {
	if !isActiveRange(state.MonthIndex, income.StartMonthIndex, income.EndMonthIndex) {
		return
	}

	// TODO: apply any inflation or up or down of the income (salary change)

	for _, ref := range income.AssetLinks {
		switch ref.TargetKind {
		case EntityAssetKind:
			for i := range state.Assets {
				if state.Assets[i].ID == ref.TargetID {
					state.Assets[i].Amount += income.Amount * float64(ref.Proportion)
				}
			}
		default:
			panic("Unimplemented: Income can not be linked to anything but asset")
		}
	}
}

func applyExpenseChange(expense *Expense, state *SimulationState) {
	if !isActiveRange(state.MonthIndex, expense.StartMonthIndex, expense.EndMonthIndex) {
		return
	}

	if expense.InflationLinked {
		expense.Amount *= 1.0 + state.Inflations[0].Value
	}

	for _, ref := range expense.AssetLinks {
		switch ref.TargetKind {
		case EntityAssetKind:
			for i := range state.Assets {
				if state.Assets[i].ID == ref.TargetID {
					state.Assets[i].Amount -= expense.Amount * float64(ref.Proportion)
				}
			}
		default:
			panic("Unimplemented: Expense can not be linked to anything but assets.")
		}
	}
}

func cleanupInflation(s *SimulationState) {
	s.Inflations = s.Inflations[1:]
}

func cleanupExpired(s *SimulationState) {
	oldAssets := s.Assets
	s.Assets = s.Assets[:0]
	for _, el := range oldAssets {
		if !isExpired(s.MonthIndex, el.EndMonthIndex) {
			s.Assets = append(s.Assets, el)
		}
	}

	oldExpenses := s.Expenses
	s.Expenses = s.Expenses[:0]

	for _, el := range oldExpenses {
		if !isExpired(s.MonthIndex, el.EndMonthIndex) {
			s.Expenses = append(s.Expenses, el)
		}
	}

	oldIncomes := s.Incomes
	s.Incomes = s.Incomes[:0]

	for _, el := range oldIncomes {
		if !isExpired(s.MonthIndex, el.EndMonthIndex) {
			s.Incomes = append(s.Incomes, el)
		}
	}

	oldLiabilities := s.Liabilities
	s.Liabilities = s.Liabilities[:0]

	for _, el := range oldLiabilities {
		if !isExpired(s.MonthIndex, el.EndMonthIndex) {
			s.Liabilities = append(s.Liabilities, el)
		}
	}
}

func simulate(state *SimulationState, month int) {
	// Trigger events based on the condition
	// for _, e := range sim.Events {
	// TODO: apply events
	// }

	// Apply inflation to the expenses
	//  - Do I apply inflation to the assets that change as well? Like a car depreciating?
	// Apply asset changes
	// Apply income changes

	for i := range state.Assets {
		applyAssetChange(&state.Assets[i], state)
	}

	for i := range state.Incomes {
		applyIncomeChange(&state.Incomes[i], state)
	}

	for i := range state.Expenses {
		applyExpenseChange(&state.Expenses[i], state)
	}

	// Create snaphost from state
	// Append snapshot to the history
	state.History = append(state.History, state.snapshot())

	// Cleanup
	// -> If we have expired items in asset, income, expense, liability remove it
	// TODO: Remove the inflation from the list!
	cleanupInflation(state)
	cleanupExpired(state)
	state.MonthIndex += 1
}

func SimulateFinancialLife(sim SimulationState, totalMonths int) []MonthlySnapshot {
	start := time.Now()
	log.Info().Msg("SIM: start")

	for month := 0; month < totalMonths; month++ {
		simulate(&sim, month)
	}

	log.Info().Msg(fmt.Sprintf("SIM: end in %s", time.Since(start)))

	return sim.History
}
