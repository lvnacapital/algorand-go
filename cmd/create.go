package cmd

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/spf13/cobra"
)

var (
	addressCmd = &cobra.Command{
		Use:   "create-address",
		Short: "Creates a wallet, and generates an account within that wallet",
		Long:  ``,
		RunE:  generateAddress,
	}

	walletCmd = &cobra.Command{
		Use:   "create-wallet",
		Short: "Creates a wallet, and generates an account within that wallet",
		Long:  ``,
		RunE:  createWallet,
	}
)

func init() {
	includeCreateFlags(walletCmd)
	includeCreateFlags(addressCmd)
}

func includeCreateFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&WalletName, "wallet", "w", "", "Wallet name")
	ccmd.Flags().StringVarP(&WalletPassword, "password", "p", "", "Wallet password")
	ccmd.Flags().BoolVarP(&generate, "generate", "g", false, "Generate a key for wallet")
}

func createWallet(ccmd *cobra.Command, args []string) error {
	if err := CheckWallet(); err != nil {
		return fmt.Errorf("Wallet name '%s' is problematic - %v", WalletName, err)
	}

	// Create the example wallet, if it doesn't already exist
	cwResponse, err := kmdClient.CreateWallet(WalletName, WalletPassword, kmd.DefaultWalletDriver, types.MasterDerivationKey{})
	if err != nil {
		return fmt.Errorf("Error creating wallet - %v", err)
	}

	// We need the wallet ID in order to get a wallet handle, so we can add accounts
	walletID := cwResponse.Wallet.ID
	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, walletID)

	if generate {
		walletHandle, err := getWalletHandle(walletID)
		if err != nil {
			return err
		}

		if err := genKey(walletHandle); err != nil {
			return err
		}
	}

	return nil
}

func generateAddress(ccmd *cobra.Command, args []string) error {
	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	walletHandle, err := GetWallet()
	if err != nil {
		return err
	}

	// Generate a new address from the wallet handle
	if err := genKey(walletHandle); err != nil {
		return err
	}

	return nil
}
