package cmd

import (
	"fmt"
	"os"

	"github.com/lvnacapital/algorand/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	getCmd = &cobra.Command{
		Use:   "status",
		Short: "Gets a transaction from the blockchain",
		Long:  ``,

		RunE: get,
	}

	findCmd = &cobra.Command{
		Use:   "status",
		Short: "Finds a transaction on the blockchain",
		Long:  ``,

		RunE: find,
	}
)

func init() {
	includeSearchFlags(findCmd)
	includeSearchFlags(getCmd)
}

func includeSearchFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&txID, "transaction", "t", "", "Specify the transaction to find")
	ccmd.Flags().StringVarP(&fromAddr, "address", "f", "", "Account address to send the money from (required)")
	ccmd.Flags().Uint64Var(&firstRound, "firstvalid", 0, "The first round where the transaction may be committed to the ledger (currently ignored)")
	ccmd.Flags().Uint64Var(&lastRound, "lastvalid", 0, "The last round where the transaction may be committed to the ledger (currently ignored)")
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
func get(ccmd *cobra.Command, args []string) error {
	util.ClearScreen()
	fmt.Println("\nFind Transaction Using Address and Transaction ID")
	fmt.Println("-------------------------------------------------")

	term := terminal.NewTerminal(os.Stdin, "")
	if fromAddr == "" {
		for {
			fmt.Print("\nEnter the sender address to look for: ")
			from, err := term.ReadLine()
			if err != nil {
				return fmt.Errorf("Error getting sender address: %s", err)
			}
			fromAddr = string(from)
			if util.IsValidAddress(fromAddr) {
				break
			}
			fmt.Print("Malformed address. Please try again.\n")
		}
	} else {
		if !util.IsValidAddress(fromAddr) {
			return fmt.Errorf("Malformed address: %s", fromAddr)
		}
	}

	if txID == "" {
		for {
			fmt.Print("\nEnter the transaction ID to look for: tx-")
			from, err := term.ReadLine()
			if err != nil {
				return fmt.Errorf("Error getting transaction ID: %s", err)
			}
			fromAddr = string(from)
			if util.IsValidAddress(fromAddr) {
				break
			}
			fmt.Print("Malformed address. Please try again.\n")
		}
	} else {
		if !util.IsValidAddress(fromAddr) {
			return fmt.Errorf("Malformed address: %s", fromAddr)
		}
	}

	tx, err := algodClient.TransactionInformation(fromAddr, txID)
	if err != nil {
		return fmt.Errorf("Transaction not found: %s", err)
	}
	fmt.Printf("Transaction: %+v", tx)

	// Reading the Note Field of a Transaction
	// Up to 1kb of arbitrary data can be stored in any
	// transaction. This data can be stored and read from the
	// transaction's note field. If this data was encoded using the
	// SDK's `encodeObj' function, then it can be decoded using the
	// `decodeObj' function.
	// const encodednote = JSON.stringify(algosdk.decodeObj(tx.note), undefined, 4);
	// fmt.Print(`Decoded: ${encodednote}`);

	// // Iterate across all transactions for a given address and
	// // round range.
	// const params = await algodClient.getTransactionParams();
	// const begin = 343496;
	// const end = params.lastRound;
	// const txts = await algodClient.transactionByAddress(addr, begin, end);
	// if (typeof txts !== 'undefined') {
	//     const lastTransaction = txts.transactions[txts.transactions.length - 1];
	//     fmt.Print(`Transaction: ${JSON.stringify(lastTransaction)}`);
	// }

	return nil
}

func find(ccmd *cobra.Command, args []string) error {
	util.ClearScreen()
	fmt.Println("\nFind Transaction Using Transaction ID")
	fmt.Println("-------------------------------------")
	// if (!txId) {
	//     txId = readlineSync.question('\nEnter the transaction ID to look for: tx-');
	// }

	// (async () => {
	//     const params = await algodClient.getTransactionParams();
	//     const start = params.lastRound;
	//     const end = 0;
	//     mainloop: for (let i = start; i > end; i--) {
	//         const block = await algodClient.block(i);
	//         // fmt.Print("Number of Transactions in " + i + ": " + block.txns.transactions.length);
	//         if (typeof block.txns.transactions === 'undefined') {
	//             continue;
	//         }
	//         const txcn = block.txns.transactions.length;

	//         for (let j = 0; j < txcn - 1; j++) {
	//             // fmt.Print("Transaction " + block.txns.transactions[j].tx);
	//             if (block.txns.transactions[j].tx === txId) {
	//                 const textedJson = JSON.stringify(block.txns.transactions[j], undefined, 4);
	//                 fmt.Print(`Transaction: ${textedJson}`);
	//                 if (
	//                     undefined !== block.txns.transactions[j].note
	//                     && block.txns.transactions[j].note.length
	//                 ) {
	//                     const encodednote = JSON.stringify(
	//                         algosdk.decodeObj(block.txns.transactions[j].note),
	//                         undefined,
	//                         4,
	//                     );
	//                     fmt.Print(`Decoded: ${encodednote}`);
	//                 }
	//                 break mainloop;
	//             }
	//         }
	//     }
	return nil
}
