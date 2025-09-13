package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

// BinanceTicker defines the structure for the incoming WebSocket message
type BinanceTicker struct {
	EventType string      `json:"e"`
	EventTime int64       `json:"c"` // This was incorrectly mapped to EventTime
	Symbol    string      `json:"s"`
	LastPrice json.Number `json:"E"` // This was incorrectly mapped to LastPrice
	Volume    json.Number `json:"v"` // Use json.Number for volume as well
}

func main() {
	// Use a channel to listen for an interrupt signal (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// The WebSocket URL for Binance's BTC/USDT ticker stream
	// See: https://binance-docs.github.io/apidocs/spot/en/#individual-symbol-ticker-streams
	url := "wss://stream.binance.com:9443/ws/btcusdt@ticker"

	log.Printf("Connecting to %s", url)

	// Dial the WebSocket server with a context
	c, _, err := websocket.DefaultDialer.DialContext(context.Background(), url, nil)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Start a goroutine to read messages from the WebSocket
	go func() {
		defer close(done)
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			if messageType != websocket.TextMessage {
				continue
			}

			var ticker BinanceTicker
			if err := json.Unmarshal(message, &ticker); err != nil {
				log.Println("unmarshal error:", err)
				continue // Continue to the next message
			}

			log.Printf("Symbol: %s, Price: %s, Volume: %s", ticker.Symbol, ticker.LastPrice, ticker.Volume)
		}
	}()

	log.Println("Successfully connected. Streaming live data... Press Ctrl+C to exit.")

	// Block until an interrupt is received or the connection is closed
	select {
	case <-done:
		log.Println("WebSocket connection closed.")
	case <-interrupt:
		log.Println("Interrupt received. Closing connection.")
		// Cleanly close the connection by sending a close message
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close error:", err)
		}
	}
}
