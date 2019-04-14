package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lvnacapital/algorand/cmd"
)

func TestBackup(t *testing.T) {
	if os.Getenv("CI") == "true" && !kmdAvailable {
		// No Algorand node available
		return
	}
	got, err := executeCommand(cmd.AlgorandCmd, "backup", "-w", walletName, "-p", walletPassword)
	if got != "" {
		expected := fmt.Sprintf("Private Key Mnemonic: %s", mnemonic)
		if got != expected {
			t.Errorf("Unexpected output - %v", got)
		}
	}
	if err != nil {
		t.Errorf("Unexpected error - %v", err)
	}
}
