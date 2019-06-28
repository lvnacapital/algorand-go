package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/lvnacapital/algorand-go/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Gets a transaction from the blockchain",
		Long:  ``,
		RunE:  getAddr,
	}

	getAllCmd = &cobra.Command{
		Use:   "get-all",
		Short: "Gets all transaction for an address on the blockchain",
		Long:  ``,
		RunE:  getAllAddr,
	}

	findCmd = &cobra.Command{
		Use:   "find",
		Short: "Finds a transaction on the blockchain",
		Long:  ``,
		RunE:  find,
	}
)

func init() {
	includeSearchFlags(getCmd)
	includeSearchFlags(getAllCmd)
	includeSearchFlags(findCmd)
}

func includeSearchFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&txID, "transaction", "t", "", "Specify the transaction to find")
	ccmd.Flags().StringVarP(&addr, "address", "a", "", "Account address to use")
	ccmd.Flags().Uint64Var(&firstRound, "firstvalid", 0, "The first round where the transaction may be found")
	ccmd.Flags().Uint64Var(&lastRound, "lastvalid", 0, "The last round where the transaction may be found")
}

func getParams() (txParams models.TransactionParams, err error) {
	// Get the suggested transaction parameters
	txParams, err = algodClient.SuggestedParams()
	if err != nil {
		err = fmt.Errorf("Error getting suggested Tx params: %s", err)
	}

	return txParams, err
}

// e.g. 2VXBXLOZSLA5EXPYD3P2SS5ODNUDTMOWTIQPLEU2SZB2Z563IWIXMQKJKI
func resolveAddress() error {
	if addr == "" {
		term := terminal.NewTerminal(os.Stdin, "")
		for {
			fmt.Print("\nEnter the address: ")
			a, err := term.ReadLine()
			if err != nil {
				return fmt.Errorf("Error getting address: %s", err)
			}
			addr = string(a)
			if util.IsValidAddress(addr) {
				break
			}
			fmt.Print("Malformed address. Please try again.\n")
		}
	} else {
		if !util.IsValidAddress(addr) {
			return fmt.Errorf("Malformed address: %s", addr)
		}
		fmt.Print("\n")
	}

	return nil
}

// e.g. A6R7R6EL2I4QJRHBSRLE2B4AQ3N74MKRWQZARYCXQOR742HC3NGQ
// e.g. Y7XIQGCRU6IRLLCKSZJVTMQSH6HSDMDVW3WQFVOMNNBNFL5BO6KQ
func resolveTxID() error {
	if txID == "" {
		term := terminal.NewTerminal(os.Stdin, "")
		fmt.Print("\nEnter the transaction ID: tx-")
		txid, err := term.ReadLine()
		if err != nil {
			return fmt.Errorf("Error getting transaction ID: %s", err)
		}
		txID = string(txid)
	} else {
		fmt.Print("\n")
	}

	return nil
}

// Reading the Note Field of a Transaction
// Up to 1kb of arbitrary data can be stored in any
// transaction. This data can be stored and read from the
// transaction's note field. If this data was encoded using the
// SDK's `Encode' function, then it can be decoded using the
// `Decode' function.
func readNote(note []byte) {
	// fmt.Printf("Note size: %d\n", len(transaction.Note))
	var m interface{}
	err := msgpack.Decode(note, &m)
	if err != nil {
		fmt.Printf("Cannot decode note - %s\n", err)
	}
	fmt.Printf("Decoded type: %T\n", m)
	fmt.Printf("Decoded byte: %v\n", m)
	if fmt.Sprintf("%T", m) == "[]uint8" {
		fmt.Printf("Decoded text: %s\n", m)
	} else if fmt.Sprintf("%T", m) == "map[interface {}]interface {}" {
		var m map[interface{}]string
		err := msgpack.Decode(note, &m)
		if err != nil {
			fmt.Printf("Cannot decode note - %s\n", err)
		}
		fmt.Print("Decoded text:\n")
		for k, v := range m {
			fmt.Printf("\t%v: %v\n", k, v)
		}
	} else {
		fmt.Println("Cannot recognize encoding")
	}
}

// Locating a Transaction
// Once a transaction is submitted and finalized into a block, it can be
// found later using several methods:
//  - If the node's 'Archival' property is not set to 'true,' only a
//    limited number of local blocks on the node are allowed to be
//    searched. If the 'Archival' property is set to 'true,' the entire
//    blockchain will be available for searching.
//  - Use the account's address with the transaction ID and call the
//    `algod' client's transactionInformation function to find a
//    specific transaction.
func getAddr(ccmd *cobra.Command, args []string) error {
	util.ClearScreen()
	fmt.Println("\nFind Transaction Using Address and Transaction ID")
	fmt.Print("-------------------------------------------------")

	resolveAddress()
	resolveTxID()

	tx, err := algodClient.TransactionInformation(fromAddr, txID)
	if err != nil {
		return fmt.Errorf("Transaction not found - %s", err)
	}
	fmt.Printf("Transaction: %+v", tx)

	return nil
}

// Iterate across all transactions for a given address and round range
func getAllAddr(ccmd *cobra.Command, args []string) error {
	util.ClearScreen()
	fmt.Println("\nList Transaction Using Address")
	fmt.Print("------------------------------")

	resolveAddress()

	txParams, err := getParams()
	if err != nil {
		return err
	}
	begin := uint64(343496)
	end := txParams.LastRound
	fmt.Printf("%s %d %d", fromAddr, begin, end)
	txts, err := algodClient.TransactionsByAddr(fromAddr, begin, end)
	if err != nil {
		return fmt.Errorf("Error finding transactions in this range - %s", err)
	}

	if len(txts.Transactions) > 0 {
		latestTransaction := txts.Transactions[len(txts.Transactions)-1]
		fmt.Printf("Latest transaction: %+v", latestTransaction)
	}

	return nil
}

func find(ccmd *cobra.Command, args []string) error {
	util.ClearScreen()
	fmt.Println("\nFind Transaction Using Transaction ID")
	fmt.Print("-------------------------------------")

	resolveTxID()

	txParams, err := getParams()
	if err != nil {
		return err
	}
	start := txParams.LastRound
	end := uint64(0)

mainLoop:
	for i := start; i >= end && i != 0; i-- {
		block, err := algodClient.Block(i)
		if err != nil {
			return fmt.Errorf("Retrieving block %d - %s", i, err)
		}
		// fmt.Printf("Number of transactions in block %d: %d", i, len(block.Transactions.Transactions))
		if !(len(block.Transactions.Transactions) > 0) {
			// fmt.Printf("No transactions in block: %d\n", i)
			continue
		}

		for _, transaction := range block.Transactions.Transactions {
			if transaction.TxID == txID {
				fmt.Printf("Found transaction in block: %d\n", i)
				transactionJSON, err := json.MarshalIndent(transaction, "", "  ")
				if err != nil {
					fmt.Printf("Cannot marshall block data - %s\n", err)
				}
				fmt.Printf("Transaction: %s\n", transactionJSON)
				if len(transaction.Note) > 0 {
					readNote(transaction.Note)
				}
				break mainLoop
			}
		}
	}

	return nil
}
