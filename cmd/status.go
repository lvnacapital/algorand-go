package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/lvnacapital/algorand-go/util"
	"github.com/spf13/cobra"
)

var (
	// Alias for status
	healthCmd = &cobra.Command{
		Hidden: true,
		Use:    "health",
		Short:  "Display the status of the Algorand node",
		Long:   ``,
		RunE:   status,
	}

	// Status command
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Display the status of the Algorand node",
		Long:  ``,
		RunE:  status,
	}
)

func status(ccmd *cobra.Command, args []string) (err error) {
	util.ClearScreen()
	// fmt.Print("`algod' Status\n--------------\n")
	// fmt.Printf("algod: %T, kmd: %T\n", algodClient, kmdClient)

	// Get algod status
	nodeStatus, err := algodClient.Status()
	if err != nil {
		return fmt.Errorf("Error getting algod status: %s", err)
	} else if os.Getenv("GOTEST") == "true" {
		ccmd.Print("Node status retrieved successfully.")
		return
	}

	fmt.Printf("Last committed block: %d\n", nodeStatus.LastRound)
	fmt.Printf("Time since last block: %.1fs\n", time.Duration.Seconds(time.Duration(nodeStatus.TimeSinceLastRound)))
	fmt.Printf("Sync Time: %.1fs\n", time.Duration.Seconds(time.Duration(nodeStatus.CatchupTime)))
	fmt.Printf("Last consensus protocol: %s\n", nodeStatus.LastVersion)
	fmt.Printf("Next consensus protocol: %s\n", nodeStatus.NextVersion)
	fmt.Printf("Round for next consensus protocol: %d\n", nodeStatus.NextVersionRound)
	fmt.Printf("Next consensus protocol supported: %t\n", nodeStatus.NextVersionSupported)
	if os.Getenv("GOTEST") == "true" {
		ccmd.Print("Node status retrieval successful.")
	}

	return
}
