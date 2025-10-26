package storage

import (
	"context"
	"github.com/pieter/GO/cryptotradingbot-starter/internal/domain"
)

// Repository defines the interface for persisting and retrieving data.
type Repository interface {
	SaveTicker(ctx context.Context, ticker domain.Ticker) error
	Close() error
}
