package exchange

import (
	"context"
	"github.com/pieter/GO/cryptotradingbot-starter/internal/domain"
)

// Streamer defines the interface for connecting to and receiving data from a live data stream.
type Streamer interface {
	Stream(ctx context.Context, symbol string) (<-chan domain.Ticker, <-chan error)
}

