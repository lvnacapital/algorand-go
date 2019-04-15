package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lvnacapital/algorand/util"
	"github.com/spf13/cobra"
)

var (
	blockCmd = &cobra.Command{
		Use:   "block",
		Short: "Display the information of a block/round",
		Long:  ``,
		RunE:  block,
	}
)

func init() {
	includeBlockFlags(blockCmd)
}

func includeBlockFlags(ccmd *cobra.Command) {
	ccmd.Flags().Uint64VarP(&blockNumber, "block", "b", uint64(0), "Block/round number to retrieve data for")
}

func block(ccmd *cobra.Command, args []string) error {
	util.ClearScreen()

	if err := getBlock(); err != nil {
		return err
	}

	blockRes, err := algodClient.Block(blockNumber)
	if err != nil {
		return fmt.Errorf("Error getting block - %s", err)
	}

	// Print the block information
	fmt.Printf("\n-----------------Block Information-------------------\n")
	blockJSON, err := json.MarshalIndent(blockRes, "", "    ")
	if err != nil {
		return fmt.Errorf("Cannot marshall block data - %s", err)
	}
	fmt.Printf("%s\n", blockJSON)
	if os.Getenv("GOTEST") == "true" {
		ccmd.Print("Block retrieved successfully.")
	}

	return nil
}
