package exchange

import (
	"context"
	"math/big"
	"time"

	"github.com/dorpsen/cryptotradingbot-starter/internal/domain"
)

// MockStreamer emits synthetic tickers at a fixed interval. Useful for local
// testing without a live exchange connection.
type MockStreamer struct {
	Interval time.Duration
	Start    float64
}

// NewMockStreamer returns a mock streamer with a default interval when zero.
func NewMockStreamer(interval time.Duration, start float64) *MockStreamer {
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	if start <= 0 {
		start = 100.0
	}
	return &MockStreamer{Interval: interval, Start: start}
}

func (m *MockStreamer) Stream(ctx context.Context, symbol string) (<-chan domain.Ticker, <-chan error) {
	tick := make(chan domain.Ticker, 8)
	errc := make(chan error, 1)

	go func() {
		defer close(tick)
		defer close(errc)

		price := big.NewFloat(m.Start)
		volume := big.NewFloat(1.0)
		t := time.NewTicker(m.Interval)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case now := <-t.C:
				// simple pseudo-random walk
				price.Add(price, big.NewFloat(0.01))
				volume.Add(volume, big.NewFloat(0.1))

				tick <- domain.Ticker{
					EventType: "mock",
					EventTime: now.UnixMilli(),
					Symbol:    symbol,
					LastPrice: domain.BigString{price},
					Volume:    domain.BigString{volume},
					OpenTime:  now.UnixMilli() - 60000,
					CloseTime: now.UnixMilli(),
					Count:     1,
				}
			}
		}
	}()

	return tick, errc
}
