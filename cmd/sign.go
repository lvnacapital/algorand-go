package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/lvnacapital/algorand/util"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

var (
	signCmd = &cobra.Command{
		Use:   "send",
		Short: "Sign and submit a transaction",
		Long:  ``,
		RunE:  sign,
	}
)

func init() {
	includeSignFlags(signCmd)
}

func includeSignFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&WalletName, "wallet", "w", "", "Set the wallet to be used for the selected operation")
	ccmd.Flags().StringVarP(&WalletPassword, "password", "p", "", "The wallet's password")
	ccmd.Flags().StringVarP(&fromAddr, "from", "f", "", "Account address to send the money from (required)")
	ccmd.Flags().StringVarP(&toAddr, "to", "t", "", "Address to send to money to (required)")
	ccmd.Flags().StringVarP(&noteText, "note", "n", "", "Note text")
	ccmd.Flags().Uint64Var(&fee, "fee", 0, "The transaction fee (automatically determined by default)")
	ccmd.Flags().Uint64VarP(&amount, "amount", "a", 0, "The filename to save the raw data to (required)")
	ccmd.Flags().Uint64Var(&firstRound, "firstvalid", 0, "The first round where the transaction may be committed to the ledger (currently ignored)")
	ccmd.Flags().Uint64Var(&lastRound, "lastvalid", 0, "The last round where the transaction may be committed to the ledger (currently ignored)")
}

// Generate a new address from the wallet handle
func generateAddrs() error {
	// gen1Response, err := kmdClient.GenerateKey(walletHandle)
	// if err != nil {
	// 	return fmt.Errorf("Error generating key: %s", err)
	// }
	// fmt.Printf("Generated address 1 %s.\n", gen1Response.Address)
	// fromAddr := gen1Response.Address

	// gen2Response, err := kmdClient.GenerateKey(walletHandle)
	// if err != nil {
	// 	return fmt.Errorf("Error generating key: %s", err)
	// }
	// fmt.Printf("Generated address 2 %s.\n", gen2Response.Address)
	// toAddr := gen2Response.Address

	return nil
}

func getFromAddr(walletHandle string) error {
	keysList, err := kmdClient.ListKeys(walletHandle)
	if err != nil {
		return fmt.Errorf("Error listing addresses - %s", err)
	} else if len(keysList.Addresses) <= 0 {
		return fmt.Errorf("No addresses available")
	}

	fmt.Printf("\nHave %d address(es) in '%s':\n", len(keysList.Addresses), WalletName)
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
		term := terminal.NewTerminal(os.Stdin, "")
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

	return nil
}

func getToAddr() error {
	if toAddr == "" {
		term := terminal.NewTerminal(os.Stdin, "")
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

	return nil
}

func getAmount() error {
	if amount == 0 {
		term := terminal.NewTerminal(os.Stdin, "")
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

	return nil
}

func getNote() error {
	if noteText == "" {
		term := terminal.NewTerminal(os.Stdin, "")
		fmt.Print("\nSpecify some note text (optional): ")
		n, err := term.ReadLine()
		if err != nil {
			return fmt.Errorf("Error getting note: %s", err)
		}
		noteText = string(n)
	}

	return nil
}

func makeTransaction() (tx *types.Transaction, err error) {
	// Get the suggested transaction parameters
	txParams, err := algodClient.SuggestedParams()
	if err != nil {
		return nil, fmt.Errorf("Error getting suggested tx params - %s", err)
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
	txP, err := transaction.MakePaymentTxn(fromAddr, toAddr, fee, amount, firstRound, lastRound, note, closeRemainderTo, genID)
	if err != nil {
		return nil, fmt.Errorf("Error creating transaction: %s", err)
	}
	tx = &txP
	fmt.Printf("\nMade transaction: %+v\n", txP)

	return
}

func signTransaction(walletHandle string, tx *types.Transaction) ([]byte, error) {
	keyRes, err := kmdClient.ExportKey(walletHandle, WalletPassword, fromAddr)
	if err != nil {
		return nil, fmt.Errorf("Error extracting secret key: %s", err)
	}
	privateKey := keyRes.PrivateKey

	// Sign the transaction (using library)
	_, stx, err := crypto.SignTransaction(privateKey, *tx)
	if err != nil {
		return nil, fmt.Errorf("Failed to sign transaction using library - %s", err)
	}
	fmt.Printf("\nMade signed transaction using library: %x\n", stx)

	// Sign the transaction (using `kmd')
	kmdStx, err := kmdClient.SignTransaction(walletHandle, WalletPassword, *tx)
	if err != nil {
		return nil, fmt.Errorf("Failed to sign transaction with `kmd' - %s", err)
	}
	fmt.Printf("\n`kmd' made signed transaction with bytes: %x\n", kmdStx.SignedTransaction)

	if bytes.Equal(kmdStx.SignedTransaction, stx) {
		fmt.Println("\nSigned transactions match")
	} else {
		return nil, fmt.Errorf("\nSigned transactions don't match")
	}

	return stx, nil
}

// Broadcast the transaction to the network
func sendTransaction(stx []byte) error {
	send, err := algodClient.SendRawTransaction(stx) // or 'kmdStx.SignedTransaction'
	if err != nil {
		return fmt.Errorf("Failed to send transaction: %s", err)
	}
	fmt.Printf("\nSent transaction with ID: tx-%s\n", send.TxID)

	return nil
}

// Signing and Submitting a Transaction
func sign(ccmd *cobra.Command, args []string) (err error) {
	walletHandle, err := GetWallet()
	if err != nil {
		return
	}

	if err = getFromAddr(walletHandle); err != nil {
		return
	}

	if err = getToAddr(); err != nil {
		return
	}

	if err = getAmount(); err != nil {
		return
	}

	if err = getNote(); err != nil {
		return
	}

	tx, err := makeTransaction()
	if err != nil {
		return
	}

	stx, err := signTransaction(walletHandle, tx)
	if err != nil {
		return
	}

	if err = sendTransaction(stx); err != nil {
		return
	}

	return nil
}
