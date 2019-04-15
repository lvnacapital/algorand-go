package cmd_test

import (
	"os"
	"testing"

	"github.com/lvnacapital/algorand/cmd"
)

func TestBlock(t *testing.T) {
	blocksGood := []string{"1", "440000"}
	blocksBad := []string{"1000000"}

	if os.Getenv("CI") == "true" && !algodAvailable {
		// No Algorand node available
		return
	} else {
		blocksGood = []string{"440000"} // not archival node
	}

	for _, block := range blocksGood {
		got, err := executeCommand(cmd.AlgorandCmd, "block", "-b", block)
		if got != "" {
			expected := "Block retrieved successfully."
			if got != expected {
				t.Errorf("Unexpected output - %v", got)
			}
		}
		if err != nil {
			t.Errorf("Unexpected error - %v", err)
		}
	}

	for _, block := range blocksBad {
		_, err := executeCommand(cmd.AlgorandCmd, "block", "-b", block)
		if err == nil {
			t.Errorf("Unexpected success - %v", err)
		}
	}
}
