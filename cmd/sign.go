package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/transaction"

	"github.com/lvnacapital/algorand/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Signing and submitting a transaction",
		Long:  ``,
		RunE:  sign,
	}
)

func init() {
	includeSignFlags(signCmd)
}

func includeSignFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&walletName, "wallet", "w", "", "Set the wallet to be used for the selected operation")
	ccmd.Flags().StringVarP(&walletPassword, "password", "p", "", "The wallet's password")
	ccmd.Flags().StringVarP(&fromAddr, "from", "f", "", "Account address to send the money from (required)")
	ccmd.Flags().StringVarP(&toAddr, "to", "t", "", "Address to send to money to (required)")
	ccmd.Flags().StringVarP(&noteText, "note", "n", "", "Note text")
	ccmd.Flags().Uint64Var(&fee, "fee", 0, "The transaction fee (automatically determined by default)")
	ccmd.Flags().Uint64VarP(&amount, "amount", "a", 0, "The filename to save the raw data to (required)")
	ccmd.Flags().Uint64Var(&firstRound, "firstvalid", 0, "The first round where the transaction may be committed to the ledger (currently ignored)")
	ccmd.Flags().Uint64Var(&lastRound, "lastvalid", 0, "The last round where the transaction may be committed to the ledger (currently ignored)")
}

// Signing and Submitting a Transaction
func sign(ccmd *cobra.Command, args []string) error {
	term := terminal.NewTerminal(os.Stdin, "")

	// Get the list of wallets
	walletsList, err := kmdClient.ListWallets()
	if err != nil {
		return fmt.Errorf("Error listing wallets: %s", err)
	} else if len(walletsList.Wallets) <= 0 {
		return fmt.Errorf("no wallets available")
	}

	var walletID string
	var walletName string
	fmt.Printf("\nGot %d wallet(s):\n", len(walletsList.Wallets))
	for i, wallet := range walletsList.Wallets {
		if walletName != "" { // Find our wallet name in the list
			if wallet.Name == walletName {
				fmt.Printf("Found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
				walletID = wallet.ID
				break
			}
		} else { // List wallets for selection
			fmt.Printf("[%d] Name: %s\tID: %s\n", i+1, wallet.Name, wallet.ID)
		}
	}
	if walletID == "" {
		for {
			if len(walletsList.Wallets) == 1 {
				fmt.Printf("Select wallet [%s]: ", "1")
			} else {
				fmt.Printf("Select wallet [%s%d]: ", "1-", len(walletsList.Wallets))
			}
			walletNum, err := term.ReadLine()
			if err != nil {
				return fmt.Errorf("Error getting wallet number: %s", err)
			}
			i, err := strconv.Atoi(string(walletNum))
			if err != nil || i > len(walletsList.Wallets) || i <= 0 {
				fmt.Print("Invalid wallet number. Please try again.\n")
				continue
			}
			walletID = walletsList.Wallets[i-1].ID
			walletName = walletsList.Wallets[i-1].Name
			break
		}
	}
	// fmt.Printf("Picked wallet %s.\n", walletID)

	fmt.Printf("Please type in the password for '%s': ", walletName)
	walletPassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("\nError getting password: %s", err)
	}
	fmt.Print("\n")

	// Get a wallet handle
	initRes, err := kmdClient.InitWalletHandle(walletID, string(walletPassword))
	if err != nil {
		return fmt.Errorf("\nError initializing wallet handle: %s", err)
	}
	walletHandle := initRes.WalletHandleToken

	keysList, err := kmdClient.ListKeys(walletHandle)
	if err != nil {
		return fmt.Errorf("Error listing addresses: %s", err)
	} else if len(keysList.Addresses) <= 0 {
		return fmt.Errorf("no addresses available")
	}

	fmt.Printf("\nGot %d address(es):\n", len(walletsList.Wallets))
	for i, address := range keysList.Addresses {
		if fromAddr != "" { // Find address in the list
			if address == fromAddr {
				fmt.Printf("Found address '%s'.\n", address)
				break
			}
		} else { // List addresses for selection
			fmt.Printf("[%d] %s\n", i+1, address)
		}
	}
	if fromAddr == "" {
		for {
			if len(keysList.Addresses) == 1 {
				fmt.Printf("Pick the account address to send from [%s]: ", "1")
			} else {
				fmt.Printf("Pick the account address to send from [%s%d]: ", "1-", len(keysList.Addresses))
			}
			addressNum, err := term.ReadLine()
			if err != nil {
				return fmt.Errorf("Error getting address number: %s", err)
			}
			i, err := strconv.Atoi(string(addressNum))
			if err != nil || i > len(keysList.Addresses) || i <= 0 {
				fmt.Print("Invalid address number. Please try again.\n")
				continue
			}
			fromAddr = keysList.Addresses[i-1]
			break
		}
	}

	keyRes, err := kmdClient.ExportKey(walletHandle, string(walletPassword), fromAddr)
	if err != nil {
		return fmt.Errorf("Error extracting secret key: %s", err)
	}
	privateKey := keyRes.PrivateKey

	// // Generate a new address from the wallet handle
	// gen1Response, err := kmdClient.GenerateKey(walletHandle)
	// if err != nil {
	// 	return fmt.Errorf("Error generating key: %s", err)
	// }
	// fmt.Printf("Generated address 1 %s.\n", gen1Response.Address)
	// fromAddr := gen1Response.Address

	// // Generate a new address from the wallet handle
	// gen2Response, err := kmdClient.GenerateKey(walletHandle)
	// if err != nil {
	// 	return fmt.Errorf("Error generating key: %s", err)
	// }
	// fmt.Printf("Generated address 2 %s.\n", gen2Response.Address)
	// toAddr := gen2Response.Address

	if toAddr == "" {
		for {
			fmt.Print("\nSpecify the account address to send to: ")
			to, err := term.ReadLine()
			if err != nil {
				return fmt.Errorf("Error getting 'to' address: %s", err)
			}
			toAddr = string(to)
			if util.IsValidAddress(toAddr) {
				break
			}
			fmt.Print("Malformed address. Please try again.\n")
		}
	}

	if amount == 0 {
		for {
			fmt.Print("\nSpecify the amount to be transferred: ")
			a, err := term.ReadLine()
			if err != nil {
				return fmt.Errorf("Error getting amount: %s", err)
			}
			amount, err = strconv.ParseUint(string(a), 10, 64)
			if err == nil && amount > 0 {
				break
			}
			fmt.Print("Invalid amount. Please try again.\n")
		}
	}

	if noteText == "" {
		fmt.Print("\nSpecify some note text (optional): ")
		n, err := term.ReadLine()
		if err != nil {
			return fmt.Errorf("Error getting note: %s", err)
		}
		noteText = string(n)
	}

	// Get the suggested transaction parameters
	txParams, err := algodClient.SuggestedParams()
	if err != nil {
		return fmt.Errorf("Error getting suggested tx params: %s", err)
	}

	// Make transaction
	if fee == 0 {
		fee = txParams.Fee
	}
	if firstRound == 0 {
		firstRound = txParams.LastRound
	}
	if lastRound == 0 {
		lastRound = txParams.LastRound + 2
	}
	if firstRound > lastRound {
		firstRound = txParams.LastRound
		lastRound = txParams.LastRound + 2
	}
	note := msgpack.Encode(noteText)
	closeRemainderTo := ""
	genID := txParams.GenesisID
	tx, err := transaction.MakePaymentTxn(fromAddr, toAddr, fee, amount, firstRound, lastRound, note, closeRemainderTo, genID)
	if err != nil {
		return fmt.Errorf("Error creating transaction: %s", err)
	}
	fmt.Printf("Made transaction: %+v\n", tx)

	// Sign the transaction (using library)
	txid, stx, err := crypto.SignTransaction(privateKey, tx)
	if err != nil {
		return fmt.Errorf("Failed to sign transaction: %s", err)
	}
	fmt.Printf("Made signed transaction with TxID %s: %x\n", txid, stx)

	// Sign the transaction (using `kmd')
	kmdStx, err := kmdClient.SignTransaction(walletHandle, string(walletPassword), tx)
	if err != nil {
		return fmt.Errorf("Failed to sign transaction with kmd: %s", err)
	}
	fmt.Printf("`kmd' made signed transaction with bytes: %x\n", kmdStx.SignedTransaction)

	if bytes.Equal(kmdStx.SignedTransaction, stx) {
		fmt.Println("Signed transactions match!")
	} else {
		fmt.Println("Signed transactions don't match!")
	}

	// Broadcast the transaction to the network
	sendResponse, err := algodClient.SendRawTransaction(stx) // or 'kmdStx.SignedTransaction'
	if err != nil {
		return fmt.Errorf("Failed to send transaction: %s", err)
	}

	fmt.Printf("Transaction ID: %s\n", sendResponse.TxID)

	return nil
}
