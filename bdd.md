Feature: Real-time Cryptocurrency Price Monitoring
  As a user, I want to monitor the live price of a cryptocurrency so that I can stay informed about market movements.
 
  Scenario: User starts monitoring a specific cryptocurrency pair
    Given the user wants to track the price of "BTC/USDT"
    When the application is started
    Then the user sees a continuous stream of price updates for "BTC/USDT"

Feature: Historical Data Persistence
  As a system, I need to store price data so that it can be used for later analysis and backtesting.

  Scenario: A new price tick is received from the data stream
    Given the application is monitoring a cryptocurrency price stream
    When a new price update is received
    Then the price update is saved to the historical data store

Feature: Application Lifecycle
  As a user, I want the application to start and stop cleanly.

  Scenario: User gracefully shuts down the application
    Given the application is streaming live data
    When the user sends a shutdown signal
    Then the connection to the data source is closed cleanly
    And the application exits without errors

  Scenario: The connection to the data source is lost
    Given the application is streaming live data
    When the connection to the data source is interrupted
    Then the application attempts to reconnect automatically

Feature: Chart Display
  As a user, I want to see a chart with CoinPair data like a graph of the active timeframe.

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
    And the active trading strategy uses a "50-period Simple Moving Average"
    When the chart is rendered
    Then a line representing the "50-period Simple Moving Average" is displayed on the price chart

Feature: Trading Opportunity
  As a user, I want to see a box where trading opportunities show up in a list for a trading pair and timeframe.

  Scenario: A new trading opportunity is identified
    Given the bot is running a strategy on "ETH/USDT" on the "1 hour" timeframe
    When the strategy's conditions for a "buy" signal are met
    Then a new entry appears in the "Trading Opportunities" list
    And the entry details the pair "ETH/USDT", the signal "Buy", and the price

Feature: Link to trading opportunity
  As a user, I need a link to a trading opportunity so it will change the GUI setup to show the data and the interface for a trade for that opportunity.

  Scenario: User navigates from an opportunity to the trade interface
    Given the "Trading Opportunities" list shows a "Buy" opportunity for "ETH/USDT"
    When the user clicks on that opportunity
    Then the main chart view changes to show "ETH/USDT"
    And the trade interface is pre-filled with details for a "Buy" order on "ETH/USDT"

Feature: Interface for trading.
  As a user, I need an interface for a trade, a box like the one on TradingView.

  Scenario: User places a market order
    Given the trade interface is open for "BTC/USDT"
    When the user enters an amount of "0.01"
    And the user clicks the "Buy Market" button
    Then a market buy order for "0.01 BTC" is sent to the exchange

Feature: Settings
  As a user, I need to be able to change some settings.

  Scenario: User changes the trading pair to monitor
    Given the application is currently monitoring "BTC/USDT"
    When the user navigates to settings and changes the monitored pair to "ETH/USDT"
    Then the application stops streaming "BTC/USDT" data
    And starts streaming "ETH/USDT" data


=====================================

Feature: Superchart

Feature: Superchart settings

Feature: 