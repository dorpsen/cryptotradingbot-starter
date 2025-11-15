package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dorpsen/cryptotradingbot-starter/internal/domain"
	"github.com/gorilla/websocket"
)

// reconnectingStreamer will continuously attempt to (re)connect to the
// provided websocket URL and forward parsed tickers to the caller.
type reconnectingStreamer struct {
	url    string
	dialer *websocket.Dialer
	// initial backoff parameters
	minBackoff time.Duration
	maxBackoff time.Duration
}

// NewReconnectingBinanceStreamer creates a Streamer that reconnects on failure.
func NewReconnectingBinanceStreamer(url string) Streamer {
	d := websocket.DefaultDialer
	d.HandshakeTimeout = 10 * time.Second
	return &reconnectingStreamer{
		url:        url,
		dialer:     d,
		minBackoff: 500 * time.Millisecond,
		maxBackoff: 30 * time.Second,
	}
}

// binanceTicker mirrors the important parts of the Binance websocket payload
// and is used only for local unmarshalling.
type binanceRawTicker struct {
	EventType string           `json:"e"`
	EventTime int64            `json:"E"`
	Symbol    string           `json:"s"`
	LastPrice domain.BigString `json:"c"`
	Volume    domain.BigString `json:"v"`
	OpenTime  int64            `json:"O"`
	CloseTime int64            `json:"C"`
	Count     int64            `json:"n"`
}

func (s *reconnectingStreamer) Stream(ctx context.Context, symbol string) (<-chan domain.Ticker, <-chan error) {
	tickerChan := make(chan domain.Ticker, 32)
	errChan := make(chan error, 1)

	go func() {
		defer close(tickerChan)
		defer close(errChan)

		backoff := s.minBackoff

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			conn, resp, err := s.dialer.DialContext(ctx, s.url, nil)
			if err != nil {
				if resp != nil {
					log.Printf("websocket handshake failed: %s", resp.Status)
				}
				// report and backoff
				select {
				case errChan <- fmt.Errorf("dial error: %w", err):
				default:
				}
				time.Sleep(backoff)
				backoff = backoff * 2
				if backoff > s.maxBackoff {
					backoff = s.maxBackoff
				}
				continue
			}

			// Reset backoff on successful connect
			backoff = s.minBackoff

			// Listen loop
			conn.SetReadLimit(65536)
			for {
				// check for ctx cancel
				select {
				case <-ctx.Done():
					conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					conn.Close()
					return
				default:
				}

				conn.SetReadDeadline(time.Now().Add(8 * time.Second))
				_, msg, err := conn.ReadMessage()
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						// Continue to allow checking ctx
						continue
					}
					// On other errors, break to reconnect
					select {
					case errChan <- fmt.Errorf("read error: %w", err):
					default:
					}
					conn.Close()
					break
				}

				var raw binanceRawTicker
				if err := json.Unmarshal(msg, &raw); err != nil {
					// If a single message fails to parse, log and continue
					log.Printf("warning: could not unmarshal message: %v", err)
					continue
				}

				tickerChan <- domain.Ticker{
					EventType: raw.EventType,
					EventTime: raw.EventTime,
					Symbol:    raw.Symbol,
					LastPrice: raw.LastPrice,
					Volume:    raw.Volume,
					OpenTime:  raw.OpenTime,
					CloseTime: raw.CloseTime,
					Count:     raw.Count,
				}
			}

			// small pause before trying to reconnect
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
	}()

	return tickerChan, errChan
}
