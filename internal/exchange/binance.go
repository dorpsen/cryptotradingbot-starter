package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pieter/GO/cryptotradingbot-starter/internal/domain"
)

// binanceTicker represents the raw data structure from the Binance API.
type binanceTicker struct {
	EventType string        `json:"e"`
	EventTime int64         `json:"E"`
	Symbol    string        `json:"s"`
	LastPrice domain.BigString `json:"c"`
	Volume    domain.BigString `json:"v"`
	OpenTime  int64         `json:"O"`
	CloseTime int64         `json:"C"`
	Count     int64         `json:"n"`
}

// toDomain converts a Binance-specific ticker to the application's generic domain.Ticker.
func (bt binanceTicker) toDomain() domain.Ticker {
	return domain.Ticker{
		EventType: bt.EventType,
		EventTime: bt.EventTime,
		Symbol:    bt.Symbol,
		LastPrice: bt.LastPrice,
		Volume:    bt.Volume,
		OpenTime:  bt.OpenTime,
		CloseTime: bt.CloseTime,
		Count:     bt.Count,
	}
}

// BinanceStreamer implements the Streamer interface for the Binance exchange.
type BinanceStreamer struct {
	conn *websocket.Conn
}

// NewBinanceStreamer creates a new streamer connected to Binance.
func NewBinanceStreamer(ctx context.Context, url string) (*BinanceStreamer, error) {
	log.Printf("Connecting to %s", url)

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	c, resp, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		if resp != nil {
			log.Printf("WebSocket handshake failed with status: %s", resp.Status)
		}
		return nil, fmt.Errorf("failed to connect to binance: %w", err)
	}
	if resp != nil {
		log.Printf("WebSocket connected with status: %s", resp.Status)
	}
	return &BinanceStreamer{conn: c}, nil
}

// Stream starts listening to the websocket and sends tickers to a channel.
func (s *BinanceStreamer) Stream(ctx context.Context, symbol string) (<-chan domain.Ticker, <-chan error) {
	tickerChan := make(chan domain.Ticker, 10)
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in websocket read: %v", r)
			}
			close(tickerChan)
			close(errChan)
			s.conn.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Context cancelled, closing websocket")
				s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			default:
			}

			s.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			_, message, err := s.conn.ReadMessage()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // It's a read timeout, just loop again to check context.
				}
				errChan <- err // Report other errors.
				return
			}

			var rawTicker binanceTicker
			if err := json.Unmarshal(message, &rawTicker); err != nil {
				log.Printf("Warning: could not unmarshal message: %v", err)
				continue
			}

			// Convert to domain object before sending
			tickerChan <- rawTicker.toDomain()
		}
	}()

	return tickerChan, errChan
}
