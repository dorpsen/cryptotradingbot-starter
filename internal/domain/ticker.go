package domain

import (
	"encoding/json"
	"math/big"
)

// Ticker represents a generic 24-hour ticker update.
// This is the primary data structure used within the application's core logic.
type Ticker struct {
	EventType string
	EventTime int64
	Symbol    string
	LastPrice BigString
	Volume    BigString
	OpenTime  int64
	CloseTime int64
	Count     int64
}

// BigString is a custom type for handling high-precision numbers from JSON strings.
type BigString struct {
	*big.Float
}

func (b *BigString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return &json.UnmarshalTypeError{}
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
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
