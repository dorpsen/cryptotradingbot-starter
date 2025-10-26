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

Feature: Indicator Visualization
  As a user, I want to see technical indicators on the chart to help with my analysis.

  Background:
    Given a price chart for "BTC/USDT" is displayed

  Scenario Outline: Overlaying indicators on the main price chart
    And the active trading strategy uses the "<indicator_name>"
    When the chart is rendered
    Then a visual representation of the "<indicator_name>" is displayed on the price chart

    Examples:
      | indicator_name                                |
      | 50-period Simple Moving Average (MA)          |
      | 21-period Exponential Moving Average (EMA)    |
      | Bollinger Bands                               |
      | Keltner Channels (KC)                         |
      | Parabolic SAR (PSAR)                          |

  Scenario Outline: Displaying indicators in a separate pane
    And the active trading strategy uses the "<indicator_name>"
    When the chart is rendered
    Then a separate pane below the price chart shows the "<indicator_name>"

    Examples:
      | indicator_name                               |
      | Moving Average Convergence Divergence (MACD) |
      | Relative Strength Index (RSI)                |
      | Stochastic Oscillator                        |
      | On-Balance Volume (OBV)                      |
      | MGHW (We gaan het meemaken)                  |

Feature: Trading Opportunity
  As a user, I want to see a list of trading opportunities so that I can quickly identify potential trades.

  Scenario: A new trading opportunity is identified
    Given the bot is running a strategy on "ETH/USDT" on the "1 hour" timeframe
    When the strategy's conditions for a "buy" signal are met
    Then a new entry appears in the "Trading Opportunities" list
    And the entry details the pair "ETH/USDT", the signal "Buy", and the price

Feature: Settings
  As a user, I need to be able to change application settings, such as the monitored trading pair.

  Scenario: User changes the trading pair to monitor
    Given the application is currently monitoring "BTC/USDT"
    When the user navigates to settings and changes the monitored pair to "ETH/USDT"
    Then the application stops streaming "BTC/USDT" data
    And starts streaming "ETH/USDT" data

Feature: Data Management
  As a system and user, I need reliable data handling for real-time monitoring and historical analysis.

  Scenario: A new price tick is received from the data stream
    Given the application is monitoring a cryptocurrency price stream
    When a new price update is received
    Then the price update is saved to the historical data store

  Scenario: User starts monitoring a specific cryptocurrency pair
    Given the user wants to track the price of "BTC/USDT"
    When the application is started
    Then the user sees a continuous stream of price updates for "BTC/USDT"

Feature: Strategy Configuration
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
