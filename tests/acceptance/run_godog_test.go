package acceptance

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/pieter/GO/cryptotradingbot-starter/tests/acceptance/harness"
	"github.com/pieter/GO/cryptotradingbot-starter/tests/acceptance/steps"
)

func TestMain(m *testing.M) {
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
