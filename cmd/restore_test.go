package cmd_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/lvnacapital/algorand/cmd"
)

func TestRestore(t *testing.T) {
	t.Parallel()
	if os.Getenv("CI") == "true" && !kmdAvailable {
		// No Algorand node available
		return
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	cmd.WalletName = walletName
	for {
		if _, err := cmd.GetWallet(); err != nil {
			walletName = fmt.Sprintf("%s-%d", walletName, r1.Intn(1000000000))
			break
		}
		time.Sleep(2 * time.Second)
	}
	got, err := executeCommand(cmd.AlgorandCmd, "restore", "-w", walletName, "-p", walletPassword, "-m", mnemonic)
	if got != "" {
		expected := "Created wallet successfully."
		if got != expected {
			t.Errorf("Unexpected output - %v", got)
		}
	}
	if err != nil {
		t.Errorf("Unexpected error - %v", err)
	}
}
