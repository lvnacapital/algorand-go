package cmd

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Create command
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Creates a wallet, and generates an account within that wallet",
		Long:  ``,

		RunE: create,
	}
)

func init() {
	includeCreateFlags(createCmd)
}

func includeCreateFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&walletName, "wallet", "w", "testwallet", "Set the wallet to be used for the selected operation")
	ccmd.Flags().StringVarP(&walletPassword, "password", "p", "testpassword", "The wallet's password")
}

func create(ccmd *cobra.Command, args []string) error {
	kmdAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port"))
	kmdToken := viper.GetString("kmd-token")

	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		return fmt.Errorf("Failed to make kmd client: %s", err)
	}
	fmt.Println("Made a kmd client")

	// Create the example wallet, if it doesn't already exist
	cwResponse, err := kmdClient.CreateWallet(walletName, walletPassword, kmd.DefaultWalletDriver, types.MasterDerivationKey{})
	if err != nil {
		return fmt.Errorf("Error creating wallet: %s", err)
	}

	// We need the wallet ID in order to get a wallet handle, so we can add accounts
	walletID := cwResponse.Wallet.ID
	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, walletID)

	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	initResponse, err := kmdClient.InitWalletHandle(walletID, walletPassword)
	if err != nil {
		return fmt.Errorf("Error initializing wallet handle: %s", err)
	}

	// Extract the wallet handle
	walletHandleToken := initResponse.WalletHandleToken

	// Generate a new address from the wallet handle
	genResponse, err := kmdClient.GenerateKey(walletHandleToken)
	if err != nil {
		return fmt.Errorf("Error generating key: %s", err)
	}
	fmt.Printf("Generated address %s\n", genResponse.Address)

	return nil
}
