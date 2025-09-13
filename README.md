# cryptotradingbot-starter

A starter project for building a cryptocurrency trading bot in Go.

This initial version connects to the Binance WebSocket API to stream live price data for a given symbol (e.g., BTC/USDT).

## Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install) (version 1.21 or later)

### Installation

1.  Clone the repository:
    ```sh
    git clone <your-repository-url>
    cd cryptotradingbot-starter
    ```

2.  Install dependencies:
    ```sh
    go mod tidy
    ```

3.  Run the application:
    ```sh
    go run main.go
    ```

You should see live price data for BTC/USDT being printed to your console. Press `Ctrl+C` to stop the stream.
