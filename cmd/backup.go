package cmd

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Backup command
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backs up a wallet, and generates an account within that wallet",
		Long:  ``,

		RunE: backup,
	}
)

func init() {
	includeBackupFlags(backupCmd)
}

func includeBackupFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&walletName, "wallet", "w", "testwallet", "Set the wallet to be used for the selected operation")
	ccmd.Flags().StringVarP(&walletPassword, "password", "p", "testpassword", "The wallet's password")
}

// You can export a master derivation key from the wallet and convert it to a mnemonic phrase in order to back up any generated addresses. This backup phrase will only allow you to recover wallet-generated keys; if you import an external key into a kmd-managed wallet, you'll need to back up that key by itself in order to recover it.
func backup(ccmd *cobra.Command, args []string) error {
	kmdAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port"))
	kmdToken := viper.GetString("kmd-token")

	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		return fmt.Errorf("Failed to make kmd client: %s", err)
	}
	fmt.Println("Made a kmd client")

	// Get the list of wallets
	listResponse, err := kmdClient.ListWallets()
	if err != nil {
		return fmt.Errorf("Error listing wallets: %s", err)
	}

	// Find our wallet name in the list
	var walletID string
	fmt.Printf("Got %d wallet(s):\n", len(listResponse.Wallets))
	for _, wallet := range listResponse.Wallets {
		fmt.Printf("ID: %s\tName: %s\n", wallet.ID, wallet.Name)
		if wallet.Name == walletName {
			fmt.Printf("Found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
			walletID = wallet.ID
		}
	}

	// Get a wallet handle
	initResponse, err := kmdClient.InitWalletHandle(walletID, walletPassword)
	if err != nil {
		return fmt.Errorf("Error initializing wallet handle: %s", err)
	}

	// Extract the wallet handle
	walletHandleToken := initResponse.WalletHandleToken

	// Get the backup phrase
	exportResponse, err := kmdClient.ExportMasterDerivationKey(walletHandleToken, walletPassword)
	if err != nil {
		return fmt.Errorf("Error exporting backup phrase: %s", err)
	}
	mdk := exportResponse.MasterDerivationKey

	// This string should be kept in a safe place and not shared
	stringToSave, err := mnemonic.FromKey(mdk[:])
	if err != nil {
		return fmt.Errorf("Error getting backup phrase: %s", err)
	}

	fmt.Printf("Backup Phrase: %s\n", stringToSave)

	return nil
}
