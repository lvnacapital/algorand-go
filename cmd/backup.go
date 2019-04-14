package cmd

import (
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Back up a wallet by stashing the private key mnemonic",
		Long:  ``,
		RunE:  backup,
	}
)

func init() {
	includeBackupFlags(backupCmd)
}

func includeBackupFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&WalletName, "wallet", "w", "", "Wallet name")
	ccmd.Flags().StringVarP(&WalletPassword, "password", "p", "", "Wallet password")
}

// Export a master derivation key from the wallet and convert it to a
// mnemonic phrase in order to back up any generated addresses.
// This backup phrase will only allow recovery of wallet-generated
// keys. When importing an external key into a kmd-managed wallet,
// it is needed to back up that key by itself in order to recover it.
func backup(ccmd *cobra.Command, args []string) (err error) {
	walletHandle, err := GetWallet()
	if err != nil {
		return
	}

	// Get the backup phrase
	exportResponse, err := kmdClient.ExportMasterDerivationKey(walletHandle, WalletPassword)
	if err != nil {
		return fmt.Errorf("Error exporting backup phrase - %s", err)
	}
	mdk := exportResponse.MasterDerivationKey

	// This string should be kept in a safe place and not shared
	stringToSave, err := mnemonic.FromKey(mdk[:])
	if err != nil {
		return fmt.Errorf("Error getting backup phrase - %s", err)
	}

	fmt.Printf("Backup Phrase: %s\n", stringToSave)
	if os.Getenv("GOTEST") == "true" {
		ccmd.Printf("Private Key Mnemonic: %s", stringToSave)
	}

	return
}
