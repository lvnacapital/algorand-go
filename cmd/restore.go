package cmd

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Restore command
	restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Backs up a wallet, and generates an account within that wallet",
		Long:  ``,

		RunE: restore,
	}
)

func init() {
	includerestoreFlags(restoreCmd)
}

func includerestoreFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&walletName, "wallet", "w", "testwallet", "Set the wallet to be used for the selected operation")
	ccmd.Flags().StringVarP(&walletPassword, "password", "p", "testpassword", "The wallet's password")
}

// To restore a wallet, convert the phrase to a key and pass it to CreateWallet. This call will fail if the wallet already exists:
func restore(ccmd *cobra.Command, args []string) error {
	kmdAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port"))
	kmdToken := viper.GetString("kmd-token")

	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		return fmt.Errorf("Failed to make kmd client: %s", err)
	}
	restorePhrase := "fire enlist diesel stamp nuclear chunk student stumble call snow flock brush example slab guide choice option recall south kangaroo hundred matrix school above zero"
	keyBytes, err := mnemonic.ToKey(restorePhrase)
	if err != nil {
		return fmt.Errorf("Failed to get key: %s", err)
	}

	var mdk types.MasterDerivationKey
	copy(mdk[:], keyBytes)
	cwResponse, err := kmdClient.CreateWallet(walletName, walletPassword, kmd.DefaultWalletDriver, mdk)
	if err != nil {
		return fmt.Errorf("Error creating wallet: %s", err)
	}
	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, cwResponse.Wallet.ID)

	return nil
}
