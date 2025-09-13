# cryptotradingbot-starter

A starter project for building a cryptocurrency trading bot in Go.

This initial version connects to the Binance WebSocket API to stream live price data for a given symbol (e.g., BTC/USDT).

## Features

*   **Complete Data Model**: Fully decodes the `24hrTicker` stream from Binance, providing access to all data fields.
*   **High-Precision Numbers**: Uses `math/big.Float` for price and volume data to prevent floating-point inaccuracies, which is crucial for financial applications.
*   **Modular Design**: WebSocket connection and streaming logic are encapsulated in a `Streamer` struct, making the code clean, reusable, and easy to test.
*   **Graceful Shutdown**: Implements context-aware handling for `Ctrl+C` interrupts, ensuring a clean closure of the WebSocket connection.
*   **Test-Driven**: Includes a unit test to verify the correctness of the data parsing logic, forming a solid foundation for future development.

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
