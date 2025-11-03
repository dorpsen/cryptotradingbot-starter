package tests

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/dorpsen/cryptotradingbot-starter/internal/domain"
	"github.com/dorpsen/cryptotradingbot-starter/internal/storage"
)

// setupTestDB creates a temporary database for the test and returns a repository and a cleanup function.
func setupTestDB(t *testing.T) (*storage.SqliteRepository, func()) {
	t.Helper()

	dbFile := "test_ticks.db"
	// Ensure any old test database file is removed before the test.
	os.Remove(dbFile)

	// Note: We now call storage.NewSqliteRepository
	repo, err := storage.NewSqliteRepository(context.Background(), dbFile)
	if err != nil {
		t.Fatalf("could not create test database: %v", err)
	}

	// The cleanup function closes the database and removes the file.
	cleanup := func() {
		repo.Close()
		os.Remove(dbFile)
	}

	return repo, cleanup
}

func TestNewSqliteRepository(t *testing.T) {
	_, cleanup := setupTestDB(t)
	defer cleanup()

	// The existence of the file and no error from setup is a sufficient test for creation.
	// We can't inspect the internal 'db' field from an external package.
	// If NewSqliteRepository failed, setupTestDB would have already failed the test.
	t.Log("NewSqliteRepository ran successfully.")
}

func TestSaveTicker(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a sample domain.Ticker object to save.
	lastPrice, _, _ := big.ParseFloat("52123.45", 10, 256, big.ToZero)
	volume, _, _ := big.ParseFloat("12345.678", 10, 256, big.ToZero)

	ticker := domain.Ticker{
		EventType: "24hrTicker",
		EventTime: 1672515782239,
		Symbol:    "BTCUSDT",
		LastPrice: domain.BigString{Float: lastPrice},
		Volume:    domain.BigString{Float: volume},
		OpenTime:  1672429382239,
		CloseTime: 1672515782239,
		Count:     1603775,
	}

	// Save the ticker to the database.
	err := repo.SaveTicker(context.Background(), ticker)
	if err != nil {
		t.Fatalf("SaveTicker failed with an unexpected error: %v", err)
	}

	// Retrieve the ticker to verify it was saved correctly.
	retrievedTicker, err := repo.GetTickerByEventTime(context.Background(), ticker.EventTime)
	if err != nil {
		t.Fatalf("GetTickerByEventTime failed with an unexpected error: %v", err)
	}
	if retrievedTicker == nil {
		t.Fatalf("expected to retrieve a ticker, but got nil")
	}

	// Compare the retrieved ticker with the original.
	// Note: Comparing big.Float requires using its Cmp method. 0 means they are equal.
	if ticker.Symbol != retrievedTicker.Symbol || ticker.EventTime != retrievedTicker.EventTime || ticker.LastPrice.Float.Cmp(retrievedTicker.LastPrice.Float) != 0 || ticker.Volume.Float.Cmp(retrievedTicker.Volume.Float) != 0 || ticker.Count != retrievedTicker.Count {
		t.Errorf("retrieved ticker does not match saved ticker.\nretrieved:  %+v\noriginal:   %+v", retrievedTicker, ticker)
	}
}
