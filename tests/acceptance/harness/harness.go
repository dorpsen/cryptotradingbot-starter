package harness

import (
	"context"
	"time"
)

// Minimal TestHarness interface to implement before using the step skeletons.
// Implement this in tests/acceptance/harness/ and provide NewTestHarness factory.

type CandleFixture struct {
	OpenTime  time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	CloseTime time.Time
}

type TickFixture struct {
	Timestamp time.Time
	Price     float64
	Volume    float64
}

type DomainEvent struct {
	Type      string
	Pair      string
	Timeframe string
	// Payload etc.
}

type ChartModel struct {
	// minimal query helpers
	candlesPresent bool
}

func (c ChartModel) HasCandles() bool { return c.candlesPresent }
func (c ChartModel) HasIndicator(name string) bool {
	return true
}

type TestHarness interface {
	SeedHistoricalData(pair string, candles []CandleFixture) error
	ConfigureMonitoredPairs(pairs []string) error
	VerifyHistoricalStoreAvailable() error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	PushTick(pair string, tick TickFixture) error
	WaitForEvent(ctx context.Context, eventType string, timeout time.Duration) (DomainEvent, error)
	QueryChart(pair, timeframe string) (ChartModel, error)
	QueryLastRenderedChart() (ChartModel, error)
	ShowChart(pair, timeframe string) error
	ConfigureStrategyAndPair(pair, timeframe string) error
	SetActiveStrategyIndicator(indicatorName string) error
	SimulateExchangeInterrupt() error
	// Additional helpers...
}

// Factory stub - implement concrete harness that wires up in-memory fakes.
func NewTestHarness() (TestHarness, error) {
	return NewInMemoryHarness()
}
