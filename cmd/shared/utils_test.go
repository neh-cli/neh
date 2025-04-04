// filename: cmd/shared/utils_test.go

package shared

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"
)

func TestExecuteWebSocketCommand(t *testing.T) {
	os.Setenv("NEH_PERSONAL_ACCESS_TOKEN", "dummy_token")

	err := ExecuteWebSocketCommand("test_command", "test_message", "clipboard_message")
	if err != nil {
		t.Errorf("ExecuteWebSocketCommand failed: %v", err)
	}
}

func TestProcessMessageInOrder(t *testing.T) {
	messagePool := &sync.Map{}
	expectedSequenceNumber := uint(1)

	// Populate the message pool with out-of-order messages
	messagePool.Store(uint(2), "second message\n")
	messagePool.Store(uint(1), "first message\n")
	messagePool.Store(uint(3), "third message\n")

	// Capture the output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	processMessageInOrder(messagePool, &expectedSequenceNumber)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	expectedOutput := "first message\nsecond message\nthird message\n"
	if buf.String() != expectedOutput {
		t.Errorf("Expected output %q, but got %q", expectedOutput, buf.String())
	}
}
