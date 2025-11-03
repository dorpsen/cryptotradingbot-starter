package storage

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/dorpsen/cryptotradingbot-starter/internal/domain"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// SqliteRepository manages the database connection and operations for SQLite.
type SqliteRepository struct {
	db *sql.DB
}

// NewSqliteRepository creates a new SqliteRepository instance and initializes the database.
func NewSqliteRepository(ctx context.Context, dbPath string) (*SqliteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	repo := &SqliteRepository{db: db}

	if err := repo.createTable(ctx); err != nil {
		return nil, err
	}

	return repo, nil
}

// createTable creates the 'ticks' table for storing ticker data.
func (s *SqliteRepository) createTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS ticks (
		event_type TEXT NOT NULL,
		event_time INTEGER NOT NULL,
		symbol TEXT NOT NULL,
		last_price TEXT NOT NULL,
		volume TEXT NOT NULL,
		open_time INTEGER NOT NULL,
		close_time INTEGER NOT NULL,
		count INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (symbol, event_time)
	);`

	_, err := s.db.ExecContext(ctx, query)
	return err
}

// SaveTicker saves a domain.Ticker object to the database. Note the change in the table schema.
func (s *SqliteRepository) SaveTicker(ctx context.Context, ticker domain.Ticker) error {
	query := `
	INSERT INTO ticks (event_type, event_time, symbol, last_price, volume, open_time, close_time, count)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	// Use full precision for price, as it's critical.
	lastPriceStr := ticker.LastPrice.Float.Text('f', -1)
	// For volume, 8 decimal places is a standard and sufficient precision.
	volumeStr := ticker.Volume.Float.Text('f', 8)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query,
		ticker.EventType, ticker.EventTime, ticker.Symbol, lastPriceStr, volumeStr,
		ticker.OpenTime, ticker.CloseTime, ticker.Count,
	)

	return err
}

// GetTickerByEventTime retrieves a ticker from the database by its event time.
func (s *SqliteRepository) GetTickerByEventTime(ctx context.Context, eventTime int64) (*domain.Ticker, error) {
	query := `SELECT event_type, event_time, symbol, last_price, volume, open_time, close_time, count FROM ticks WHERE event_time = ?`
	row := s.db.QueryRowContext(ctx, query, eventTime)

	var ticker domain.Ticker
	var lastPriceStr, volumeStr string

	err := row.Scan(
		&ticker.EventType,
		&ticker.EventTime,
		&ticker.Symbol,
		&lastPriceStr,
		&volumeStr,
		&ticker.OpenTime,
		&ticker.CloseTime,
		&ticker.Count,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found is a valid outcome, not an error.
		}
		return nil, fmt.Errorf("could not scan ticker row: %w", err)
	}

	// Convert string representations back to big.Float
	ticker.LastPrice.Float, _, err = big.ParseFloat(lastPriceStr, 10, 256, big.ToZero)
	if err != nil {
		return nil, fmt.Errorf("could not parse last_price: %w", err)
	}
	ticker.Volume.Float, _, err = big.ParseFloat(volumeStr, 10, 256, big.ToZero)
	if err != nil {
		return nil, fmt.Errorf("could not parse volume: %w", err)
	}

	return &ticker, nil
}

// Close closes the database connection.
func (s *SqliteRepository) Close() error {
	return s.db.Close()
}
