package harness

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// Concrete minimal in-memory TestHarness implementation used by acceptance tests.
// This is intentionally small and deterministic to support fast godog tests.

type inMemoryHarness struct {
	mu         sync.Mutex
	historical map[string][]CandleFixture
	monitored  []string
	events     chan DomainEvent
	started    bool
	clock      *ManualClock
	stopped    chan struct{}
}

func NewInMemoryHarness() (TestHarness, error) {
	h := &inMemoryHarness{
		historical: make(map[string][]CandleFixture),
		events:     make(chan DomainEvent, 100),
		clock:      NewManualClock(time.Now()),
		stopped:    make(chan struct{}),
	}
	return h, nil
}

// ManualClock: simple controllable clock for tests
type ManualClock struct {
	mu   sync.Mutex
	time time.Time
}

func NewManualClock(start time.Time) *ManualClock {
	return &ManualClock{time: start}
}

func (c *ManualClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.time
}

func (c *ManualClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.time = c.time.Add(d)
	c.mu.Unlock()
}

// Implementation of TestHarness interface

func (h *inMemoryHarness) SeedHistoricalData(pair string, candles []CandleFixture) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.historical[pair] = append([]CandleFixture{}, candles...)
	return nil
}

func (h *inMemoryHarness) ConfigureMonitoredPairs(pairs []string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.monitored = append([]string{}, pairs...)
	return nil
}

func (h *inMemoryHarness) VerifyHistoricalStoreAvailable() error {
	// always available for in-memory store
	return nil
}

func (h *inMemoryHarness) Start(ctx context.Context) error {
	h.mu.Lock()
	if h.started {
		h.mu.Unlock()
		return nil
	}
	h.started = true
	h.mu.Unlock()

	// Publish startup events for monitored pairs
	go func() {
		for _, p := range h.monitored {
			h.events <- DomainEvent{Type: "ExchangeConnected", Pair: p}
		}
		h.events <- DomainEvent{Type: "HistoricalStoreInitialized"}
		h.events <- DomainEvent{Type: "StreamingStarted"}
	}()

	return nil
}

func (h *inMemoryHarness) Stop(ctx context.Context) error {
	h.mu.Lock()
	if !h.started {
		h.mu.Unlock()
		return nil
	}
	h.started = false
	h.mu.Unlock()

	// Publish disconnect event
	h.events <- DomainEvent{Type: "ExchangeDisconnected"}
	close(h.stopped)
	return nil
}

func (h *inMemoryHarness) PushTick(pair string, tick TickFixture) error {
	// For tests, store tick as a trivial one-candle sequence (or append as last candle)
	h.mu.Lock()
	defer h.mu.Unlock()
	c := CandleFixture{
		OpenTime:  tick.Timestamp,
		Open:      tick.Price,
		High:      tick.Price,
		Low:       tick.Price,
		Close:     tick.Price,
		Volume:    tick.Volume,
		CloseTime: tick.Timestamp,
	}
	h.historical[pair] = append(h.historical[pair], c)
	// Publish an event that a new tick/candle arrived
	h.events <- DomainEvent{Type: "TickReceived", Pair: pair}
	return nil
}

func (h *inMemoryHarness) WaitForEvent(ctx context.Context, eventType string, timeout time.Duration) (DomainEvent, error) {
	deadline := time.After(timeout)
	for {
		select {
		case <-ctx.Done():
			return DomainEvent{}, ctx.Err()
		case <-deadline:
			return DomainEvent{}, errors.New("timeout waiting for event")
		case ev := <-h.events:
			// match by exact type or type:pair pattern
			if ev.Type == eventType || ev.Type == eventType+":"+ev.Pair {
				return ev, nil
			}
		}
	}
}

func (h *inMemoryHarness) QueryChart(pair, timeframe string) (ChartModel, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	candles := h.historical[pair]
	if len(candles) == 0 {
		return ChartModel{candlesPresent: false}, nil
	}
	return ChartModel{candlesPresent: len(candles) > 0}, nil
}

func (h *inMemoryHarness) QueryLastRenderedChart() (ChartModel, error) {
	// Return a basic chart model for assertions
	return ChartModel{candlesPresent: true}, nil
}

func (h *inMemoryHarness) ShowChart(pair, timeframe string) error {
	// publish chart updated
	h.events <- DomainEvent{Type: "ChartUpdated", Pair: pair}
	return nil
}

func (h *inMemoryHarness) ConfigureStrategyAndPair(pair, timeframe string) error {
	// record config; for simplicity, just ensure the pair is monitored
	return h.ConfigureMonitoredPairs([]string{pair})
}

func (h *inMemoryHarness) SetActiveStrategyIndicator(indicatorName string) error {
	// no-op for minimal harness
	h.events <- DomainEvent{Type: "StrategyIndicatorSet"}
	return nil
}

func (h *inMemoryHarness) SimulateExchangeInterrupt() error {
	// publish disconnect then a reconnect attempt
	h.events <- DomainEvent{Type: "ExchangeDisconnected"}
	// simulate reconnect attempt soon after
	go func() {
		time.Sleep(100 * time.Millisecond)
		h.events <- DomainEvent{Type: "ReconnectAttempt"}
	}()
	return nil
}

// Helper to allow JSON docstring parsing into CandleFixture
func ParseCandlesJSON(data []byte) ([]CandleFixture, error) {
	var raw []struct {
		OpenTime  string  `json:"open_time"`
		Open      float64 `json:"open"`
		High      float64 `json:"high"`
		Low       float64 `json:"low"`
		Close     float64 `json:"close"`
		Volume    float64 `json:"volume"`
		CloseTime string  `json:"close_time"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	out := make([]CandleFixture, 0, len(raw))
	for _, r := range raw {
		ot, _ := time.Parse(time.RFC3339, r.OpenTime)
		ct, _ := time.Parse(time.RFC3339, r.CloseTime)
		out = append(out, CandleFixture{
			OpenTime:  ot,
			Open:      r.Open,
			High:      r.High,
			Low:       r.Low,
			Close:     r.Close,
			Volume:    r.Volume,
			CloseTime: ct,
		})
	}
	return out, nil
}
