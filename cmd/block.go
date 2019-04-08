package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/lvnacapital/algorand/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Block command
	blockCmd = &cobra.Command{
		Use:   "status",
		Short: "Display information of a block",
		Long:  ``,

		RunE: block,
	}
)

func init() {
	includeBlockFlags(blockCmd)
}

func includeBlockFlags(ccmd *cobra.Command) {
	ccmd.Flags().Uint64VarP(&blockNumber, "block", "b", 0, "The block number to retrieve data for")
}

func block(ccmd *cobra.Command, args []string) error {
	algodAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("algod-port"))
	kmdAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port"))
	algodToken := viper.GetString("algod-token")
	kmdToken := viper.GetString("kmd-token")

	if runtime.GOOS != "windows" {
		cli.ClearScreen()
	}

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

	fmt.Printf("algod: %T, kmd: %T\n", algodClient, kmdClient)

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
