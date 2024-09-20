// filename: cmd/decache_test.go

package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunDecacheCmd(t *testing.T) {
	// Prepare for capturing stdout
	var outputBuffer bytes.Buffer
	rootCmd.SetOut(&outputBuffer)

	// Execute the decache command
	decacheCmd.Run(&cobra.Command{}, []string{})

	// Check the output
	expectedOutput := "" // Since we don't expect any specific output, it should be empty
	if outputBuffer.String() != expectedOutput {
		t.Errorf("Expected output %q, but got %q", expectedOutput, outputBuffer.String())
	}
}
