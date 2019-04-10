package cmd

import (
	"fmt"

	"github.com/lvnacapital/algorand/util"
	"github.com/spf13/cobra"
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
	util.ClearScreen()
	fmt.Printf("\nalgod: %T, kmd: %T\n", algodClient, kmdClient)

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
