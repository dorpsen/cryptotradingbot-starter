package app

import (
	"context"
	"log"

	"github.com/dorpsen/cryptotradingbot-starter/internal/exchange"
	"github.com/dorpsen/cryptotradingbot-starter/internal/storage"
)

// Application holds the core components and orchestrates the application's logic.
type Application struct {
	streamer exchange.Streamer
	repo     storage.Repository
}

// New creates a new Application.
func New(streamer exchange.Streamer, repo storage.Repository) *Application {
	return &Application{
		streamer: streamer,
		repo:     repo,
	}
}

// Run starts the main application loop.
func (a *Application) Run(ctx context.Context, symbol string) error {
	log.Println("Application starting...")

	tickerChan, errChan := a.streamer.Stream(ctx, symbol)
	log.Println("Successfully connected. Live data is being streamed...")

	for {
		select {
		case ticker, ok := <-tickerChan:
			if !ok {
				log.Println("Ticker stream has stopped.")
				return nil
			}
			log.Printf("Symbol: %s, Price: %s", ticker.Symbol, ticker.LastPrice.Float.Text('f', 2))

			if err := a.repo.SaveTicker(ctx, ticker); err != nil {
				log.Printf("Error saving ticker: %v", err)
			}
		case err := <-errChan:
			log.Printf("Stream error: %v", err)
			return err
		case <-ctx.Done():
			log.Println("Application shutting down.")
			return nil
		}
	}
}
