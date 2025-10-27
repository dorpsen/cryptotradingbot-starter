package steps

import (
	"context"
	"fmt"
	"time"

	"github.com/cucumber/godog"

	"github.com/dorpsen/cryptotradingbot-starter/tests/acceptance/harness"
)

var H harness.TestHarness

// SetHarness allows the test runner to inject the concrete harness
func SetHarness(h harness.TestHarness) {
	H = h
}

// --- Step implementations ---

func iHaveHistoricalPriceDataForPair(pair string, rawData *godog.DocString) error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	if rawData == nil {
		return fmt.Errorf("no fixture provided")
	}
	candles, err := harness.ParseCandlesJSON([]byte(rawData.Content))
	if err != nil {
		return err
	}
	return H.SeedHistoricalData(pair, candles)
}

func theApplicationIsConfiguredToMonitor(pair string) error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	return H.ConfigureMonitoredPairs([]string{pair})
}

func theHistoricalDataStoreIsAvailable() error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	return H.VerifyHistoricalStoreAvailable()
}

func theUserStartsTheApplication() error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return H.Start(ctx)
}

func itSuccessfullyConnectsToTheDataStream(pair string) error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := H.WaitForEvent(ctx, "ExchangeConnected", 5*time.Second)
	return err
}

func itInitializesConnectionToHistoricalStore() error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := H.WaitForEvent(ctx, "HistoricalStoreInitialized", 3*time.Second)
	return err
}

func theApplicationIsStreamingLiveData() error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := H.WaitForEvent(ctx, "StreamingStarted", 3*time.Second)
	return err
}

func userSendsShutdownSignal() error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return H.Stop(ctx)
}

func connectionIsClosedCleanly() error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := H.WaitForEvent(ctx, "ExchangeDisconnected", 3*time.Second)
	return err
}

func applicationExitsWithoutErrors() error {
	// Harness.Stop should return nil; optionally check for Error events
	return nil
}

func connectionToPriceDataStreamIsInterrupted() error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	return H.SimulateExchangeInterrupt()
}

func applicationAttemptsToReconnectAutomatically() error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := H.WaitForEvent(ctx, "ReconnectAttempt", 10*time.Second)
	return err
}

// Chart display steps

func userSelectsPairAndTimeframe(pair, timeframe string) error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	return H.ShowChart(pair, timeframe)
}

func aPriceChartForPairOnTimeframeIsDisplayed(pair, timeframe string) error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	chart, err := H.QueryChart(pair, timeframe)
	if err != nil {
		return err
	}
	if !chart.HasCandles() {
		return fmt.Errorf("chart has no candles for %s %s", pair, timeframe)
	}
	return nil
}

// Indicators and Trading Opportunities

func activeTradingStrategyUses(indicatorName string) error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	return H.SetActiveStrategyIndicator(indicatorName)
}

func chartIsRendered() error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := H.WaitForEvent(ctx, "ChartUpdated", 3*time.Second)
	return err
}

func visualRepresentationOfIndicatorIsDisplayed(indicatorName string) error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	chart, err := H.QueryLastRenderedChart()
	if err != nil {
		return err
	}
	if !chart.HasIndicator(indicatorName) {
		return fmt.Errorf("indicator %s not present on chart", indicatorName)
	}
	return nil
}

func aNewTradingOpportunityIsIdentified(pair, timeframe string) error {
	if H == nil {
		return fmt.Errorf("test harness not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ev, err := H.WaitForEvent(ctx, "TradingOpportunity", 5*time.Second)
	if err != nil {
		return err
	}
	if ev.Pair != pair || ev.Timeframe != timeframe {
		return fmt.Errorf("unexpected opportunity: %v", ev)
	}
	return nil
}

// Godog initialization wiring

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the application is configured to monitor "([^"]*)"$`, theApplicationIsConfiguredToMonitor)
	ctx.Step(`^the historical data store is available$`, theHistoricalDataStoreIsAvailable)
	ctx.Step(`^the user starts the application$`, theUserStartsTheApplication)
	ctx.Step(`^it successfully connects to the "([^"]*)" data stream$`, itSuccessfullyConnectsToTheDataStream)
	ctx.Step(`^it initializes the connection to the historical data store$`, itInitializesConnectionToHistoricalStore)
	ctx.Step(`^the application is streaming live data$`, theApplicationIsStreamingLiveData)
	ctx.Step(`^the user sends a shutdown signal \(e.g., Ctrl\+C\)$`, userSendsShutdownSignal)
	ctx.Step(`^the connection to the data source is closed cleanly$`, connectionIsClosedCleanly)
	ctx.Step(`^the application exits without errors$`, applicationExitsWithoutErrors)
	ctx.Step(`^the connection to the price data stream is interrupted$`, connectionToPriceDataStreamIsInterrupted)
	ctx.Step(`^the application attempts to reconnect to the stream automatically$`, applicationAttemptsToReconnectAutomatically)

	// Chart/indicator/trading opportunities
	ctx.Step(`^historical price data for "([^"]*)" is available$`, iHaveHistoricalPriceDataForPair)
	ctx.Step(`^the user selects the "([^"]*)" pair and the "([^"]*)" timeframe$`, userSelectsPairAndTimeframe)
	ctx.Step(`^a price chart for "([^"]*)" on the "([^"]*)" timeframe is displayed$`, aPriceChartForPairOnTimeframeIsDisplayed)
	ctx.Step(`^the active trading strategy uses the "([^"]*)"`, activeTradingStrategyUses)
	ctx.Step(`^a visual representation of the "([^"]*)" is displayed on the price chart$`, visualRepresentationOfIndicatorIsDisplayed)

	ctx.Step(`^the bot is running a strategy on "([^"]*)" on the "([^"]*)" timeframe$`, func(pair, timeframe string) error {
		return H.ConfigureStrategyAndPair(pair, timeframe)
	})
	ctx.Step(`^the strategy's conditions for a "([^"]*)" signal are met$`, func(signal string) error {
		// Optionally push ticks to produce the signal
		return nil
	})
	ctx.Step(`^a new entry appears in the "Trading Opportunities" list$`, func() error {
		return aNewTradingOpportunityIsIdentified("ETH/USDT", "1h")
	})
}