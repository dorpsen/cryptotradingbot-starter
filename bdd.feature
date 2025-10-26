Feature: Application Lifecycle
  As a user, I want the application to manage its operational state reliably, from startup through execution to shutdown.

  Scenario: User starts the application successfully
    Given the application is configured to monitor "BTC/USDT"
    And the historical data store is available
    When the user starts the application
    Then it successfully connects to the "BTC/USDT" data stream
    And it initializes the connection to the historical data store

  Scenario: User gracefully shuts down the application
    Given the application is streaming live data
    When the user sends a shutdown signal (e.g., Ctrl+C)
    Then the connection to the data source is closed cleanly
    And the application exits without errors

  Scenario: The connection to the price data stream is lost
    Given the application is streaming live data
    When the connection to the price data stream is interrupted
    Then the application attempts to reconnect to the stream automatically

Feature: Chart Display
  As a user, I want to see a price chart for a trading pair so that I can visually analyze market trends.

  Scenario: User views a price chart for a specific trading pair
    Given historical price data for "BTC/USDT" is available
    When the user selects the "BTC/USDT" pair and the "15 minute" timeframe
    Then a price chart for "BTC/USDT" on the "15 minute" timeframe is displayed
    And the chart updates automatically as new price data arrives

Feature: Error Handling
  As a system, I need a way to deal with errors or wrong settings for actions.

  Scenario: User enters an invalid symbol to monitor
    Given the user is in the settings panel
    When the user tries to monitor an invalid symbol like "INVALIDCOIN"
    Then the application shows an error message "Invalid symbol specified"
    And the application does not attempt to start a new data stream

Feature: Indicators Display
  As a user, I want to see a display of the indicators used by the strategy, aligned in time with the chart.

  Scenario: An indicator is overlaid on the price chart
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses a "50-period Simple Moving Average" (MA)
    When the chart is rendered
    Then a line representing the "50-period Simple Moving Average" is displayed on the price chart

  Scenario: Exponential Moving Average (EMA) is overlaid on the price chart
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses a "21-period Exponential Moving Average" (EMA)
    When the chart is rendered
    Then a line representing the "21-period EMA" is displayed on the price chart

  Scenario: Bollinger Bands are overlaid on the price chart
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses "Bollinger Bands" with a 20-period SMA and 2 standard deviations
    When the chart is rendered
    Then three lines representing the upper, middle, and lower Bollinger Bands are displayed on the price chart

  Scenario: Keltner Channels (KC) are overlaid on the price chart
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses "Keltner Channels"
    When the chart is rendered
    Then three lines representing the upper, middle, and lower Keltner Channels are displayed on the price chart

  Scenario: Parabolic SAR (PSAR) is overlaid on the price chart
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses "Parabolic SAR"
    When the chart is rendered
    Then a series of dots representing the Parabolic SAR is displayed above or below the price candles

  Scenario: MACD is displayed in a separate pane
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses the "Moving Average Convergence Divergence" (MACD) indicator
    When the chart is rendered
    Then a separate pane below the price chart shows the MACD line, signal line, and histogram

  Scenario: RSI is displayed in a separate pane
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses the "Relative Strength Index" (RSI)
    When the chart is rendered
    Then a separate pane below the price chart shows the RSI line, typically with overbought and oversold levels

  Scenario: Stochastic Oscillator is displayed in a separate pane
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses the "Stochastic Oscillator"
    When the chart is rendered
    Then a separate pane below the price chart shows the Stochastic Oscillator with its %K and %D lines

  Scenario: On-Balance Volume (OBV) is displayed in a separate pane
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses the "On-Balance Volume" (OBV) indicator
    When the chart is rendered
    Then a separate pane below the price chart shows the OBV line

  Scenario: A custom indicator (MGHW) is displayed
    Given a price chart for "BTC/USDT" is displayed
    And the active trading strategy uses the custom "MGHW (We gaan het meemaken)" indicator
    When the chart is rendered
    Then a separate pane below the price chart shows the MGHM symbols
    
Feature: Trading Opportunity
  As a user, I want to see a list of trading opportunities so that I can quickly identify potential trades.

  Scenario: A new trading opportunity is identified
    Given the bot is running a strategy on "ETH/USDT" on the "1 hour" timeframe
    When the strategy's conditions for a "buy" signal are met
    Then a new entry appears in the "Trading Opportunities" list
    And the entry details the pair "ETH/USDT", the signal "Buy", and the price

Feature: Trade Interface
  As a user, I need a comprehensive interface to view my positions, manage alerts, and execute trades.

  Background:
    Given the trade interface layout is based on the 'Trade.md' wireframe

  Scenario: Viewing current position and balance
    Given the trade interface is open for "BTC/USDT"
    And the user has an open position of "0.5 BTC"
    And the user has a "Free" balance of "10000 USDT"
    And the user has a "Reserved" balance of "0 USDT"
    When the interface is displayed
    Then the "Position Info" section shows "0.5 BTC" for "BTC/USDT"
    And the "Balance" section shows a "Free" balance of "10000 USDT"

  Scenario: User navigates from an opportunity to the trade interface
    Given the "Trading Opportunities" list shows a "Buy" opportunity for "ETH/USDT"
    When the user clicks on that opportunity
    Then the main chart view changes to show "ETH/USDT"
    And the trade interface is pre-filled with details for a "Buy" order on "ETH/USDT"

  Scenario: Placing a limit order from the 'Trade' tab
    Given the user is on the "Trade" tab of the trade interface for "BTC/USDT"
    And the "Limit" order tab is active
    When the user enters "50000" in the "Limit Price" field
    And the user enters "100" in the "Spend" field
    Then the "Receive" field is automatically calculated to show "0.002"
    And the user clicks the final "Buy BTC" button
    Then a limit buy order for "0.002 BTC" at "50000 USDT" is sent to the exchange

  Scenario: Placing a market order using the percentage slider
    Given the user is on the "Trade" tab of the trade interface for "BTC/USDT"
    And the user has a "Free" balance of "1000 USDT"
    When the user selects the "Market" order tab
    And the user moves the percentage slider to "50%"
    Then the "Spend" field shows approximately "500"
    And the user clicks the "Buy BTC" button
    Then a market buy order for "500 USDT" worth of "BTC" is sent to the exchange

  Scenario: Staging a trade using a preset from the 'Presets' tab
    Given the user is on the "Presets" tab of the trade interface
    And a preset named "Scalp 1% Profit" exists
    When the user selects the "Scalp 1% Profit" preset for "ETH/USDT"
    And the user clicks "Apply Preset"
    Then a set of pre-configured orders (e.g., entry, take-profit, stop-loss) for "ETH/USDT" is staged for execution

  Scenario: Creating a trade from an alert in the 'Alerts' tab
    Given the user is on the "Alerts" tab of the trade interface
    And a "Buy" opportunity for "ADA/USDT" is listed based on the "RSI Divergence" strategy
    When the user selects the "ADA/USDT" opportunity
    Then the trade parameters are pre-filled based on the "RSI Divergence" strategy rules
    And the user can review and confirm the trade setup

Feature: Trade Execution Resilience
  As a trader, I want the application to handle trading connection errors gracefully, so that I can be confident my orders are processed correctly or I am notified of failures.

  Scenario: System reconnects to the trading API to place an order
    Given the user is ready to place a "buy" order for "BTC/USDT"
    And the connection to the exchange's trading API is temporarily unavailable
    When the user confirms the order
    Then the application attempts to re-establish the connection to the trading API
    And upon successful reconnection, the "buy" order is placed
    And the user receives a confirmation that the order was placed successfully

  Scenario: System fails to reconnect to the trading API and notifies the user
    Given the user is ready to place a "buy" order for "BTC/USDT"
    And the connection to the exchange's trading API is down and cannot be restored
    When the user confirms the order
    Then the application attempts to re-establish the connection to the trading API
    And after failing to reconnect, it displays an error message: "Failed to connect to exchange. Order was not placed."
    And the order is not sent

Feature: Settings
  As a user, I need to be able to change application settings, such as the monitored trading pair.

  Scenario: User changes the trading pair to monitor
    Given the application is currently monitoring "BTC/USDT"
    When the user navigates to settings and changes the monitored pair to "ETH/USDT"
    Then the application stops streaming "BTC/USDT" data
    And starts streaming "ETH/USDT" data

Feature: Historical Data Persistence
  As a system, I need to store price data so that it can be used for later analysis and backtesting.

  Scenario: A new price tick is received from the data stream
    Given the application is monitoring a cryptocurrency price stream
    When a new price update is received
    Then the price update is saved to the historical data store

Feature: Real-time Cryptocurrency Price Monitoring
  As a user, I want to monitor the live price of a cryptocurrency so that I can stay informed about market movements.

  Scenario: User starts monitoring a specific cryptocurrency pair
    Given the user wants to track the price of "BTC/USDT"
    When the application is started
    Then the user sees a continuous stream of price updates for "BTC/USDT"

Feature: Trading Strategy
  As a user, I want to be able to choose and adjust the parameters of a trading strategy. 
  For the Presets Tab and on the Alerts Tab.

Feature: Test Trading strategy
  As a user, I want to be able to test a trading strategy.

  Scenario: User Starts test tradingstrategy
    Given the user wants to test the strategy "Cryptocoiners Strategie 2.0"
    When the user has chosen the right strategy 
    And Coinpair
    And period of time
    Then the user sees a graph of the profit/loss of that strategy over that periode of time 
    And a table wit the trading results.
