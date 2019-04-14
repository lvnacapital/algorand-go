package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/lvnacapital/algorand/util"
	"github.com/spf13/cobra"
)

var (
	blockCmd = &cobra.Command{
		Use:   "block",
		Short: "Display the information of a block",
		Long:  ``,
		RunE:  block,
	}
)

func init() {
	includeBlockFlags(blockCmd)
}

func includeBlockFlags(ccmd *cobra.Command) {
	ccmd.Flags().Uint64VarP(&blockNumber, "block", "b", 0, "Block number to retrieve data for")
}

func block(ccmd *cobra.Command, args []string) error {
	util.ClearScreen()

	// Get algod status
	nodeStatus, err := algodClient.Status()
	if err != nil {
		return fmt.Errorf("Error getting algod status: %s", err)
	}

	if nodeStatus.LastRound < blockNumber {
		return fmt.Errorf("Block number cannot be greater than last round: %d > %d", blockNumber, nodeStatus.LastRound)
	}
	// Print the block information
	fmt.Printf("\n-----------------Block Information-------------------\n")
	blockJSON, err := json.MarshalIndent(blockNumber, "", "\t")
	if err != nil {
		return fmt.Errorf("Cannot marshall block data: %s", err)
	}
	fmt.Printf("%s\n", blockJSON)

	return nil
}
