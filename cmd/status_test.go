package cmd_test

import (
	"testing"

	"github.com/lvnacapital/algorand-go/cmd"
)

func TestStatus(t *testing.T) {
	got, err := executeCommand(cmd.AlgorandCmd, "status")
	if got != "" {
		expected := "Node status retrieved successfully."
		if got != expected {
			t.Errorf("Unexpected output - %v", got)
		}
	}
	if err != nil {
		t.Errorf("Unexpected error - %v", err)
	}
}
