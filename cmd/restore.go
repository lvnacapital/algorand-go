package cmd

import (
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/spf13/cobra"
)

var (
	restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Backs up a wallet, and generates an account within that wallet",
		Long:  ``,
		RunE:  restore,
	}
)

func init() {
	includeRestoreFlags(restoreCmd)
}

func includeRestoreFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&WalletName, "wallet", "w", "", "Wallet name")
	ccmd.Flags().StringVarP(&WalletPassword, "password", "p", "", "Wallet password")
	ccmd.Flags().StringVarP(&WalletMnemonic, "mnemonic", "m", "", "Private key mnemonic")
}

// To restore a wallet, convert the phrase to a key and pass it to CreateWallet. This call will fail if the wallet already exists:
func restore(ccmd *cobra.Command, args []string) error {
	keyBytes, err := getPrivateKey()
	if err != nil {
		return err
	}

	var mdk types.MasterDerivationKey
	copy(mdk[:], keyBytes)
	cwResponse, err := kmdClient.CreateWallet(WalletName, WalletPassword, kmd.DefaultWalletDriver, mdk)
	if err != nil {
		return fmt.Errorf("Error creating wallet - %s", err)
	}

	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, cwResponse.Wallet.ID)
	if os.Getenv("GOTEST") == "true" {
		ccmd.Print("Created wallet successfully.")
	}

	return nil
}
