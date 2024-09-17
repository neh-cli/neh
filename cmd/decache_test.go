// filename: cmd/decache_test.go
package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunDecacheCmd(t *testing.T) {
	os.Setenv("NEH_PERSONAL_ACCESS_TOKEN", "dummy_token")

	// Mock the ExecuteWebSocketCommand function
	mockExecuteWebSocketCommand := func(command, message string, waitForResponse bool) error {
		if command != "decache" {
			t.Errorf("Expected command 'decache', got %s", command)
		}
		if message != "" {
			t.Errorf("Expected empty message, got %s", message)
		}
		return nil
	}

	// Create a new decache command
	cmd := &cobra.Command{}
	cmd.SetArgs([]string{"decache"})

	// Capture the output
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	// Execute the command
	runDecacheCmd(cmd, []string{}, mockExecuteWebSocketCommand)

	// Check the output
	if output.String() != "" {
		t.Errorf("Expected no output, got %s", output.String())
	}
}
