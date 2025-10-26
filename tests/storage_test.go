package tests

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/pieter/GO/cryptotradingbot-starter/internal/domain"
	"github.com/pieter/GO/cryptotradingbot-starter/internal/storage"
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

	// Verification from an external package is more complex.
	// For this test, we'll rely on the lack of an error as a sign of success.
	// A more advanced test could involve adding a "Get" method to the repository interface
	// to retrieve and verify the data.
	t.Log("SaveTicker ran without error.")
}
