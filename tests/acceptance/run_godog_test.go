package acceptance

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/dorpsen/cryptotradingbot-starter/tests/acceptance/harness"
	"github.com/dorpsen/cryptotradingbot-starter/tests/acceptance/steps"
)

func TestMain(m *testing.M) {
	// By default, acceptance tests (godog) are skipped unless explicitly enabled
	// via the RUN_ACCEPTANCE=1 environment variable. This keeps `go test ./...`
	// fast for local development while still allowing acceptance runs in CI.
	if os.Getenv("RUN_ACCEPTANCE") != "1" {
		os.Exit(m.Run())
	}
	// Initialize harness and inject into steps package
	h, err := harness.NewInMemoryHarness()
	if err != nil {
		panic(err)
	}
	steps.SetHarness(h)

	status := godog.TestSuite{
		Name:                "acceptance",
		ScenarioInitializer: steps.InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"features"}, // Paths are relative to this file's directory
			Randomize:     0,
			StopOnFailure: false,
			Output:        colors.Colored(os.Stdout),
		},
	}.Run()

	// Run unit tests as well and combine status
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
