package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lvnacapital/algorand/cmd"
)

func TestBackup(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		// No Algorand node connected
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
