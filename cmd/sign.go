package cmd

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Sign command
	signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Signing and submitting a transaction",
		Long:  ``,

		RunE: sign,
	}
)

func init() {
	includeSignFlags(signCmd)
}

func includeSignFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&walletName, "wallet", "w", "testwallet", "Set the wallet to be used for the selected operation")
	ccmd.Flags().StringVarP(&walletPassword, "password", "p", "testpassword", "The wallet's password")
	ccmd.Flags().StringVarP(&fromAddr, "from", "f", "", "Account address to send the money from (required)")
	ccmd.Flags().StringVarP(&toAddr, "to", "t", "testpassword", "Address to send to money to (required)")
	ccmd.Flags().StringVarP(&note, "note", "n", "", "Note text")
	ccmd.Flags().Uint64Var(&fee, "fee", 1, "The transaction fee (automatically determined by default)")
	ccmd.Flags().Uint64VarP(&amount, "amount", "a", 0, "The filename to save the raw data to (required)")
	ccmd.Flags().Uint64Var(&firstRound, "firstvalid", 0, "The first round where the transaction may be committed to the ledger (currently ignored)")
	ccmd.Flags().Uint64Var(&lastRound, "lastvalid", 0, "The last round where the transaction may be committed to the ledger (currently ignored)")
}

func sign(ccmd *cobra.Command, args []string) error {
	algodAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("algod-port"))
	kmdAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port"))
	algodToken := viper.GetString("algod-token")
	kmdToken := viper.GetString("kmd-token")

	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return fmt.Errorf("Failed to make algod client: %s", err)
	}
	fmt.Println("Made an algod client.")

	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		return fmt.Errorf("Failed to make kmd client: %s", err)
	}
	fmt.Println("Made a kmd client.")

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

	// Generate a new address from the wallet handle
	gen1Response, err := kmdClient.GenerateKey(walletHandleToken)
	if err != nil {
		return fmt.Errorf("Error generating key: %s", err)
	}
	fmt.Printf("Generated address 1 %s.\n", gen1Response.Address)
	fromAddr := gen1Response.Address

	// Generate a new address from the wallet handle
	gen2Response, err := kmdClient.GenerateKey(walletHandleToken)
	if err != nil {
		return fmt.Errorf("Error generating key: %s", err)
	}
	fmt.Printf("Generated address 2 %s.\n", gen2Response.Address)
	toAddr := gen2Response.Address

	// Get the suggested transaction parameters
	txParams, err := algodClient.SuggestedParams()
	if err != nil {
		return fmt.Errorf("Error getting suggested tx params: %s", err)
	}

	// Make transaction
	genID := txParams.GenesisID
	// MakePaymentTxn(from, to, fee, amount, firstRound, lastRound, note, closeRemainderTo, genesisID)
	tx, err := transaction.MakePaymentTxn(fromAddr, toAddr, 1, 100, 300, 400, nil, "", genID)
	if err != nil {
		return fmt.Errorf("Error creating transaction: %s", err)
	}

	// Sign the transaction
	signResponse, err := kmdClient.SignTransaction(walletHandleToken, "testpassword", tx)
	if err != nil {
		return fmt.Errorf("Failed to sign transaction with kmd: %s", err)
	}

	fmt.Printf("kmd made signed transaction with bytes: %x\n", signResponse.SignedTransaction)

	// Broadcast the transaction to the network
	// Note that this transaction will get rejected because the accounts do not have any tokens
	sendResponse, err := algodClient.SendRawTransaction(signResponse.SignedTransaction)
	if err != nil {
		return fmt.Errorf("Failed to send transaction: %s", err)
	}

	fmt.Printf("Transaction ID: %s\n", sendResponse.TxID)

	return nil
}
