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

func status(ccmd *cobra.Command, args []string) error {
	algodAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("algod-port"))
	kmdAddress := fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port"))
	algodToken := viper.GetString("algod-token")
	kmdToken := viper.GetString("kmd-token")

	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return fmt.Errorf("Failed to make algod client: %s", err)
	}
	fmt.Println("Made an algod client")

	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		return fmt.Errorf("Failed to make kmd client: %s", err)
	}
	fmt.Println("Made a kmd client")

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
