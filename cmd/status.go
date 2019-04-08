package cmd

import (
	"fmt"
	"runtime"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/lvnacapital/algorand/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Alias for status
	healthCmd = &cobra.Command{
		Hidden: true,

		Use:   "health",
		Short: "Display the status of the Algorand node",
		Long:  ``,

		RunE: status,
	}

	// Status command
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Display the status of the Algorand node",
		Long:  ``,

		RunE: status,
	}
)

// init
func init() {
	// includeShowFlags(fetchCmd)
	// includeShowFlags(getCmd)
	// includeShowFlags(showCmd)
}

func includeShowFlags(cmd *cobra.Command) {
	// cmd.Flags().StringVarP(&key, "key", "k", "", "The key to get the data by")
	// cmd.Flags().StringVarP(&file, "file", "f", "", "The filename to save the raw data to")
	// cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print more information about request")
	// cmd.Flags().BoolVarP(&insecure, "insecure", "i", insecure, "Whether or not to ignore hoarder certificate.")
}

func status(ccmd *cobra.Command, args []string) error {
	algodAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("algod-port"))
	kmdAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port"))
	algodToken := viper.GetString("algod-token")
	kmdToken := viper.GetString("kmd-token")

	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return fmt.Errorf("")
	}

	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		return fmt.Errorf("")
	}

	if runtime.GOOS != "windows" {
		cli.ClearScreen()
	}
	fmt.Printf("algod: %T, kmd: %T\n", algodClient, kmdClient)

	// Get algod status
	nodeStatus, err := algodClient.Status()
	if err != nil {
		return fmt.Errorf("Error getting algod status: %s", err)
	}

	fmt.Printf("algod last round: %d\n", nodeStatus.LastRound)
	fmt.Printf("algod time since last round: %d\n", nodeStatus.TimeSinceLastRound)
	fmt.Printf("algod catchup: %d\n", nodeStatus.CatchupTime)
	fmt.Printf("algod latest version: %s\n", nodeStatus.LastVersion)

	return nil
}
