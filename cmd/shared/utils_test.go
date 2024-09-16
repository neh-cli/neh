// filename: cmd/shared/utils_test.go

package shared

import (
	"os"
	"testing"
)

func TestExecuteWebSocketCommand(t *testing.T) {
	os.Setenv("NEH_PERSONAL_ACCESS_TOKEN", "dummy_token")

	err := ExecuteWebSocketCommand("test_command", "test_message", false)
	if err != nil {
		t.Errorf("ExecuteWebSocketCommand failed: %v", err)
	}
}
