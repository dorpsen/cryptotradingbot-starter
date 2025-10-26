package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/pieter/GO/cryptotradingbot-starter/internal/domain"

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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_time INTEGER NOT NULL,
		symbol TEXT NOT NULL,
		last_price TEXT NOT NULL,
		volume TEXT NOT NULL,
		open_time INTEGER NOT NULL,
		close_time INTEGER NOT NULL,
		trade_count INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.ExecContext(ctx, query)
	return err
}

// SaveTicker saves a domain.Ticker object to the database.
func (s *SqliteRepository) SaveTicker(ctx context.Context, ticker domain.Ticker) error {
	query := `
	INSERT INTO ticks (event_time, symbol, last_price, volume, open_time, close_time, trade_count)
	VALUES (?, ?, ?, ?, ?, ?, ?);`

	lastPriceStr := ticker.LastPrice.Float.Text('f', -1)
	volumeStr := ticker.Volume.Float.Text('f', -1)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query,
		ticker.EventTime, ticker.Symbol, lastPriceStr, volumeStr,
		ticker.OpenTime, ticker.CloseTime, ticker.Count,
	)

	return err
}

// Close closes the database connection.
func (s *SqliteRepository) Close() error {
	return s.db.Close()
}
