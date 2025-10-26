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

Feature: Error Handling
  As a system, I need a way to deal with errors or wrong settings for actions.

  Scenario: User enters an invalid symbol to monitor
    Given the user is in the settings panel
    When the user tries to monitor an invalid symbol like "INVALIDCOIN"
    Then the application shows an error message "Invalid symbol specified"
    And the application does not attempt to start a new data stream

Feature: Chart Display
  As a user, I want to see a price chart for a trading pair so that I can visually analyze market trends.

  Scenario: User views a price chart for a specific trading pair
    Given historical price data for "BTC/USDT" is available
    When the user selects the "BTC/USDT" pair and the "15 minute" timeframe
    Then a price chart for "BTC/USDT" on the "15 minute" timeframe is displayed
    And the chart updates automatically as new price data arrives

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

  Background: TradingOpportunities.wireframe describes the layout of the Trading Opportunities list.
  <settings>: Change the settings of the Trading Opportunities to be displayed
  <strategy descreption>: The description of the strategy. Exc.: Stoch & Bollinger Bands without mimal margin
  [icon signaltrend]: The type of trend the signal belongs to. Example: Bullish (green dot), Bearish (red dot)
  <timeframe>: The timeframe of the chart.
  <Exchangename>: Name of the exchange where data is coming from in abbreviation. Example: MEXC
  <kind of moment>: Description of action according to signal. Example: Possible buy moment.
  <candlestamp date>: Date of the candlestick closure. Example: 2025-10-02 (YYYY-MM-DD)
  <candlestamp time>: Time of the candlestick closure. Example: 20:58 (HH:MM)
  <candlestamp price>: Closing price of the candlestick. Example: 1.165 USDT
  <volume24h>: Volume of the last 24 hours. Example: 3814656.28 USDT
  <volume24h price change>: Price change of the last 24 hours. Example: -2.43%
  <BB%width>: Bollinger Bands width. Example: 1.4%
  <stoch %K>: Stochastic %K. Example: 0%
  <stoch %D>: Stochastic %D. Example: 16%
  [icon chart-trend]: Chart-trend icon for graphical representation of the charttrend. Example: Bullish (green line graph going up)
  <description chart-trend>: Description of the trend. Example: Bullish.
  [icon market-trend]: Market-trend icon for graphical representation of the markettrend. Example: Bearish (red line graph going down).
  <Exchange name>: Name of the exchange that provided the data for markettrend. Example: MEXC.
  <markettrend%>: The percentage the market changed. Example 11.3%
  <markettrend description>: Description of the markettrend. Example: Bearish.
  <Description WGHM indicator>: Description of the WGHM indicator. Example: None - Neutral.
  <time opportunity> |: Time of the opportunity. Example: 20:59 (HH:MM).
  
  Strategy settings
  - select coinpairs by filtering on:
    a. 24h trading volume, minimum
    b. 24h trading volume, maximum
    c. 24h change, minimum (bearish)
    d. 24h change, maximum (bullish)

  Scenario: A new trading opportunity is identified
    Given the application is running a strategy on "ETH/USDT" on the "1 hour" timeframe
    When the strategy's conditions for a "buy" signal are met
    Then a new entry appears in the "Trading Opportunities" at the bottoum of the list
    And the entry details according to the TradingOpportunities wireframe filled in.

  Scenario: A Trading Opportunity is clicked
    Given the "Trading Opportunities" list shows a "Buy" opportunity for "ETH/USDT" on the 1m timeframe
    When the user clicks on that opportunity
    And the user clicks OK on a popup
    Then the main chart view changes to show "ETH/USDT" on the 1m timeframe
    And the trade interface is pre-filled with details for a "Buy" order on "ETH/USDT"

Feature: Settings
  As a user, I need to be able to change application settings, such as the monitored trading pairs.

  Scenario: User changes the trading pair to monitor
    Given the application is currently monitoring "BTC/USDT"
    When the user navigates to settings and changes the monitored pair to include "ETH/USDT"
    And the user changes the monitored pair to exclude "BTC/USDT"
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

Feature: Markettrend

  Scenario: Application is started
    Given the application is monitoring a cryptocurrency price stream
    And no previous Markettrend is available
    And historical data is available for the lowest timeframe
    When the application is started
    Then the charttrend is calculated for the lowest timeframe
    And for any other timeframe
    And the cart trend data per timeframe is saved
    And the overall chart trend of the cryptocurrency is calculated
    
  Scenario: A new price tick is recieved from the data system
    Given the application is monitoring a cryptocurrency price stream
    And a previous chart trend is available
    When a new price update is received 
    Then the charttrend is calculated for the lowest timeframe
    And for any other timeframe which a candle is closed by that data
    And the overall chart trend of the cryptocurrency is updated or created.

Feature: Chart Trend

  Scenario: Application is started
    Given the application is monitoring a cryptocurrency price stream
    And no previous chart trend is available
    And historical data is available for the lowest timeframe
    When the application is started
    Then the charttrend is calculated for the lowest timeframe
    And for any other timeframe
    And the cart trend data per timeframe is saved
    And the overall chart trend of the cryptocurrency is calculated
    
  Scenario: A new price tick is recieved from the data system
    Given the application is monitoring a cryptocurrency price stream
    And a previous chart trend is available
    When a new price update is received 
    Then the charttrend is calculated for the lowest timeframe
    And for any other timeframe which a candle is closed by that data
    And the overall chart trend of the cryptocurrency is updated or created.

Feature: Test Markettrend  
  As a user, I want to be able to test the Marktettrend.

  Scenario: User Starts test Marktettrend
    Given the user wants to test the Markettrend
    When the user has chosen the time and period for the tes to calculate the Markettrend for
    Then the user sees a graph of markettrend

Feature: Test Chart Trend
  As a user, I want to be able to test the Chart Trend.

  Scenario: User Starts test Chart Trend
  Given the user wants to test the Chart Trend
  When the user has chosen the time and period for the test to calculate the Chart Trend for
  Then the user sees a graph of the Chart Trend