package main

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// BinanceTicker defines the structure for the incoming WebSocket message
type BinanceTicker struct {
	EventType          string    `json:"e"`
	EventTime          int64     `json:"E"`
	Symbol             string    `json:"s"`
	PriceChange        string    `json:"p"`
	PriceChangePercent string    `json:"P"`
	WeightedAvgPrice   string    `json:"w"`
	PrevClosePrice     string    `json:"x"`
	LastPrice          BigString `json:"c"` // custom type to parse string numbers
	LastQty            string    `json:"Q"`
	BidPrice           string    `json:"b"`
	BidQty             string    `json:"B"`
	AskPrice           string    `json:"a"`
	AskQty             string    `json:"A"`
	OpenPrice          string    `json:"o"`
	HighPrice          string    `json:"h"`
	LowPrice           string    `json:"l"`
	Volume             BigString `json:"v"` // custom type
	QuoteVolume        string    `json:"q"`
	OpenTime           int64     `json:"O"`
	CloseTime          int64     `json:"C"`
	FirstID            int64     `json:"F"`
	LastID             int64     `json:"L"`
	Count              int64     `json:"n"`
}

// BigString wraps big.Float and implements json.Unmarshaler for string numbers
type BigString struct {
	*big.Float
}

func (b *BigString) UnmarshalJSON(data []byte) error {
	// Explicitly check for JSON null
	if string(data) == "null" {
		return &json.UnmarshalTypeError{}
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	// Handle empty string case, which is valid in some API responses
	if s == "" {
		b.Float = big.NewFloat(0)
		return nil
	}
	f, _, err := big.ParseFloat(s, 10, 256, big.ToZero)
	if err != nil {
		return err
	}
	b.Float = f
	return nil
}

// Streamer handles the connection and message streaming from the WebSocket.
type Streamer struct {
	conn *websocket.Conn
}

// NewStreamer creates and connects a new Streamer.
func NewStreamer(ctx context.Context, url string) (*Streamer, error) {
	c, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return &Streamer{conn: c}, nil
}

// Stream starts listening for messages and sends them to a channel.
func (s *Streamer) Stream(ctx context.Context) (<-chan BinanceTicker, <-chan error) {
	tickerChan := make(chan BinanceTicker)
	errChan := make(chan error, 1)

	go func() {
		defer close(tickerChan)
		defer close(errChan)
		defer s.conn.Close()

		for {
			// Set a deadline for the next read.
			// This makes ReadMessage non-blocking and allows us to check the context.
			s.conn.SetReadDeadline(time.Now().Add(1 * time.Second))

			_, message, err := s.conn.ReadMessage()
			if err != nil {
				// If it's a timeout error, continue the loop to check the context again.
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Before continuing, check if the context was cancelled during the wait.
					if ctx.Err() != nil {
						s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
						return // Exit the goroutine to signal shutdown.
					}
					continue
				}
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					errChan <- err
				}
				return
			}

			var ticker BinanceTicker
			if err := json.Unmarshal(message, &ticker); err != nil {
				log.Printf("Warning: could not unmarshal message: %v", err)
				continue
			}

			// Send the successfully parsed ticker.
			tickerChan <- ticker
		}
	}()

	return tickerChan, errChan
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	url := "wss://stream.binance.com:9443/ws/btcusdt@ticker"
	log.Printf("Connecting to %s", url)

	streamer, err := NewStreamer(ctx, url)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	tickerChan, errChan := streamer.Stream(ctx)
	log.Println("Successfully connected. Streaming live data... Press Ctrl+C to exit.")

	for {
		select {
		case ticker := <-tickerChan:
			log.Printf("Symbol: %s, Price: %s", ticker.Symbol, ticker.LastPrice.Float.Text('f', 2))
		case err := <-errChan:
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Error: %v", err)
			}
			log.Println("Stream closed.")
			return
		}
	}
}
