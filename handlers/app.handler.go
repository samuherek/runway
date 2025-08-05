package handlers

import (
	"fmt"
	"runway/db"
	"runway/engine"
	"runway/services/notifications"
	"runway/views/app"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AppHandler struct {
	db *db.DbService
}

func NewAppHandler(db *db.DbService) *AppHandler {
	return &AppHandler{
		db: db,
	}
}

func (h *AppHandler) Home(c echo.Context) error {
	view := app.Home()
	return renderView(c, app.HomePage(view))
}

func (h *AppHandler) GetSimplePrediction(c echo.Context) error {
	return renderView(c, app.SimplePredictionPage(app.SimplePrediction()))
}

type PostSimplePrediction struct {
	InitialBalance     float64 `form:"initialBalance" validate:"required,min=1"`
	MonthlyIncome      float64 `form:"monthlyIncome" validate:"min=0"`
	MonthlyExpenses    float64 `form:"monthlyExpenses" validate:"required,min=1"`
	ExpensesConfidence float64 `form:"expensesConfidence" validate:"min=0,max=1"`
}

func (h *AppHandler) PostSimplePrediction(c echo.Context) error {
	var params PostSimplePrediction
	n := notifications.NewNotifications()

	if !isHxReq(c) {
		return notHxResponse(c, n)
	}

	if err := c.Bind(&params); err != nil {
		log.Error().Err(err).Msg("Failed input binding")
		return unexpectedErrResponse(c, n)
	}

	fmt.Printf("Params %v\n", params)
	if err := c.Validate(&params); err != nil {
		return validationResponse(c, n, intoValidationMessages(err))
	}

	if params.ExpensesConfidence == 0 {
		params.ExpensesConfidence = 0.8
	}

	input := engine.Input{
		InitialBalance:     params.InitialBalance,
		MonthlyIncome:      params.MonthlyIncome,
		MonthlyExpenses:    params.MonthlyExpenses,
		ExpensesConfidence: params.ExpensesConfidence,
		PlannedExpenses:    make([]engine.OneTimeExpense, 0),
		UnexpectedExpense: engine.UnexpectedExpense{
			Probability: 0.25,
			MaxCost:     3000.0,
			Frequency:   1,
		},
		InflationRates: []engine.InflationForecast{
			{Year: 2025, Rate: 0.05},
			{Year: 2026, Rate: 0.04},
			{Year: 2027, Rate: 0.03},
			{Year: 2028, Rate: 0.02},
		},
		MaxMonths:   120,
		Simulations: 1000,
	}

	data := engine.SimulateSimpleProjection(input)
	d := engine.ExtractMinMax(data)
	var dates []string
	var mins []float64
	var mids []float64
	var maxs []float64
	var minDist, maxDist string

	for _, item := range d {
		dates = append(dates, item.Date)
		mins = append(mins, item.Min)
		mids = append(mids, item.Mid)
		maxs = append(maxs, item.Max)

		if item.Min <= 0 && minDist == "" {
			minDist = item.Date
		}

		if item.Max <= 0 && maxDist == "" {
			maxDist = item.Date
		}
	}

	return renderView(c, app.Chart(dates, mins, mids, maxs))
}

func (h *AppHandler) GetRetireProjection(c echo.Context) error {
	return renderView(c, app.RetireProjectionPage(app.RetireProjection()))
}

type PostRetireProjection struct {
	MonthlyExpenses      float64 `form:"monthlyExpenses" validate:"required,min=1"`
	YearsUntilRetirement int     `form:"yearsUntilRetirement" validate:"required,min=1"`
	WithdrawalYears      int     `form:"withdrawalYears" validate:"required,min=1"`
	CurrentSavings       float64 `form:"currentSavings" validate:"required,min=0"`
}

func (h *AppHandler) PostRetireProjection(c echo.Context) error {
	var params PostRetireProjection
	n := notifications.NewNotifications()

	if !isHxReq(c) {
		return notHxResponse(c, n)
	}

	if err := c.Bind(&params); err != nil {
		log.Error().Err(err).Msg("Failed input binding")
		return unexpectedErrResponse(c, n)
	}

	fmt.Printf("Params %v\n", params)
	if err := c.Validate(&params); err != nil {
		return validationResponse(c, n, intoValidationMessages(err))
	}

	input := engine.RetirementInput{
		MonthlyExpenseToday:  params.MonthlyExpenses,
		YearsUntilRetirement: params.YearsUntilRetirement,
		InflationRate:        0.03,
		WithdrawalYears:      params.WithdrawalYears,
		CurrentSavings:       params.CurrentSavings,
	}

	data := engine.ProjectRetirement(input)

	fmt.Printf("Monthly expense at retirement start (in future €): %.2f\n", data.MonthlyExpenseAtStart)
	fmt.Printf("Total required funds at retirement (future €): %.2f\n", data.TotalRequiredFutureFund)
	fmt.Printf("Total required funds in today's money: %.2f\n", data.TotalRequiredPresentValueFund)
	fmt.Printf("Current savings coverage (%% of required present value): %.2f%%\n", data.CurrentCoveragePercentage)

	return renderView(c, app.RetireResult())
}
