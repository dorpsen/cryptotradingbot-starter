package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/dorpsen/cryptotradingbot-starter/internal/app"
	"github.com/dorpsen/cryptotradingbot-starter/internal/exchange"
	"github.com/dorpsen/cryptotradingbot-starter/internal/storage"
	_ "github.com/mattn/go-sqlite3" // Driver for database/sql
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// --- Wiring Layer ---
	// Create concrete implementations of our services.
	dbPath := "ticks.db"
	repo, err := storage.NewSqliteRepository(ctx, dbPath)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer repo.Close()
	log.Printf("Database successfully opened at %s", dbPath)

	// The symbol would come from config in a real app.
	symbol := "btcusdt"
	url := "wss://stream.binance.com:9443/ws/" + strings.ToLower(symbol) + "@ticker"

	streamer, err := exchange.NewBinanceStreamer(ctx, url)
	if err != nil {
		log.Fatalf("Streamer connection failed: %v", err)
	}

	// Create the main application object, injecting the dependencies.
	application := app.New(streamer, repo)

	// Run the application.
	if err := application.Run(ctx, symbol); err != nil {
		log.Fatalf("Application run failed: %v", err)
	}

	log.Println("Application finished gracefully.")
}
