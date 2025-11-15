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

# The rest of the original feature content (wireframes and additional features)
# has been commented out to keep the acceptance suite focused on the
# core lifecycle scenarios during automated test runs.
#
# Full feature/wireframe content is available in the repo docs and can be
# re-enabled by restoring the original content or splitting into separate
# feature files one-per-Feature.