package cmd

import (
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
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
	ccmd.Flags().StringVarP(&walletName, "wallet", "w", "", "Wallet name")
	ccmd.Flags().StringVarP(&walletPassword, "password", "p", "", "Wallet password")
	ccmd.Flags().StringVarP(&walletMnemonic, "mnemonic", "m", "", "Private key mnemonic")
}

func getPrivateKey() (keyBytes []byte, err error) {
	if walletMnemonic == "" {
		term := terminal.NewTerminal(os.Stdin, "")
		for {
			fmt.Print("\nEnter the wallet mnemonic: ")
			m, err := term.ReadLine()
			if err != nil {
				return nil, fmt.Errorf("Error getting mnemonic - %v", err)
			}
			walletMnemonic = string(m)
			if keyBytes, err = mnemonic.ToKey(walletMnemonic); err != nil {
				fmt.Printf("Failed to get key. Try again - %v", err)
				continue
			}
			break
		}
	} else {
		if keyBytes, err = mnemonic.ToKey(walletMnemonic); err != nil {
			return nil, fmt.Errorf("Failed to get key from -m: %v", err)
		}
	}

	return
}

// To restore a wallet, convert the phrase to a key and pass it to CreateWallet. This call will fail if the wallet already exists:
func restore(ccmd *cobra.Command, args []string) error {
	keyBytes, err := getPrivateKey()
	if err != nil {
		return err
	}

	var mdk types.MasterDerivationKey
	copy(mdk[:], keyBytes)
	cwResponse, err := kmdClient.CreateWallet(walletName, walletPassword, kmd.DefaultWalletDriver, mdk)
	if err != nil {
		return fmt.Errorf("Error creating wallet: %s", err)
	}

	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, cwResponse.Wallet.ID)
	if os.Getenv("GOTEST") == "true" {
		ccmd.Print("Created wallet successfully.")
	}

	return nil
}
