package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestUnmarshalBinanceTicker(t *testing.T) {
	// Sample JSON payload from Binance documentation for a ticker stream
	payload := []byte(`{
		"e": "24hrTicker",
		"E": 1672515782239,
		"s": "BTCUSDT",
		"p": "273.94000000",
		"P": "1.665",
		"w": "16638.62439448",
		"x": "16456.12000000",
		"c": "16730.06000000",
		"Q": "0.010",
		"b": "16729.58000000",
		"B": "1.639",
		"a": "16730.06000000",
		"A": "0.493",
		"o": "16456.12000000",
		"h": "16758.00000000",
		"l": "16438.22000000",
		"v": "202438.852",
		"q": "3368269838.131",
		"O": 1672429382239,
		"C": 1672515782239,
		"F": 2779322214,
		"L": 2780925988,
		"n": 1603775
	}`)

	var ticker BinanceTicker
	err := json.Unmarshal(payload, &ticker)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if ticker.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", ticker.Symbol)
	}

	if expectedPrice := "16730.06"; ticker.LastPrice.Float.Text('f', 2) != expectedPrice {
		t.Errorf("expected price %s, got %s", expectedPrice, ticker.LastPrice.Float.Text('f', 2))
	}
}

func TestBigString_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    []byte
		expected    string
		expectError bool
	}{
		{"valid float string", []byte(`"123.45"`), "123.45", false},
		{"valid integer string", []byte(`"789"`), "789", false},
		{"empty string becomes zero", []byte(`""`), "0", false},
		{"invalid string", []byte(`"not-a-number"`), "", true},
		{"json null", []byte(`null`), "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bs BigString
			err := json.Unmarshal(tc.jsonData, &bs)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if bs.Float.Text('f', -1) != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, bs.Float.Text('f', -1))
			}
		})
	}
}

var upgrader = websocket.Upgrader{}

func mockBinanceServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Sample ticker message to send to the client
	tickerMessage := []byte(`{"e":"24hrTicker","s":"BTCUSDT","c":"50000.00"}`)

	// Send one message to test streaming
	if err := conn.WriteMessage(websocket.TextMessage, tickerMessage); err != nil {
		return
	}

	// Keep the connection open until the client closes it
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break // Exit loop if client disconnects
		}
	}
}

func TestStreamer(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(mockBinanceServer))
	defer server.Close()

	// Convert http:// to ws://
	url := "ws" + strings.TrimPrefix(server.URL, "http")

	t.Run("streams data and shuts down gracefully", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Test NewStreamer
		streamer, err := NewStreamer(ctx, url)
		if err != nil {
			t.Fatalf("NewStreamer failed: %v", err)
		}

		tickerChan, errChan := streamer.Stream(ctx)

		select {
		case ticker := <-tickerChan:
			if ticker.Symbol != "BTCUSDT" {
				t.Errorf("expected symbol BTCUSDT, got %s", ticker.Symbol)
			}
			if expectedPrice := "50000.00"; ticker.LastPrice.Float.Text('f', 2) != expectedPrice {
				t.Errorf("expected price %s, got %s", expectedPrice, ticker.LastPrice.Float.Text('f', 2))
			}
		case err := <-errChan:
			t.Fatalf("received unexpected error from stream: %v", err)
		case <-ctx.Done():
			t.Fatal("test timed out before receiving a ticker")
		}

		// Cancel the context to test graceful shutdown
		cancel()

		// Check that the error channel reports a normal closure
		select {
		case err := <-errChan:
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				t.Errorf("expected normal closure, but got: %v", err)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("test timed out waiting for stream to close")
		}
	})
}
