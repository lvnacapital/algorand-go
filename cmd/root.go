package cmd

import (
	"fmt"
	"os"

	"github.com/lvnacapital/algorand/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
)

var (
	// Root variables
	config      string
	showVersion bool

	algodClient algod.Client
	kmdClient   kmd.Client

	// Subcommand variables
	walletName     string
	walletPassword string
	blockNumber    uint64
	fromAddr       string
	toAddr         string
	noteText       string
	fee            uint64
	amount         uint64
	firstRound     uint64
	lastRound      uint64
	addr           string
	txID           string

	// Linker variables
	version string
	commit  string

	// AlgorandCmd ...
	AlgorandCmd = &cobra.Command{
		Use:               "algorand",
		Short:             "algorand - node explorer",
		Long:              ``,
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: allPreFlight,
		PreRunE:           rootPreFlight,
		RunE:              startAlgorand,
	}
)

func readConfig() {
	if config != "" {
		// Use config file passed in the flag.
		viper.SetConfigFile(config)
	} else {
		// // Find home directory.
		// homedir "github.com/mitchellh/go-homedir"
		// dir, err := homedir.Dir()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in "home" or "pwd" directory with name "config.yml".
		viper.AddConfigPath(dir)
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Cannot read config:", err)
		os.Exit(1)
	}
}

func allPreFlight(ccmd *cobra.Command, args []string) (err error) {
	// // Convert the log level.
	// logLvl := lumber.LvlInt(viper.GetString("log-level"))

	// // Configure the logger.
	// lumber.Prefix("[algorand]")
	// lumber.Level(logLvl)

	nodeConfig := util.Node{
		AlgodAddress: fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("algod-port")),
		KmdAddress:   fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port")),
		AlgodToken:   viper.GetString("algod-token"),
		KmdToken:     viper.GetString("kmd-token"),
	}

	algodClient, kmdClient, err = util.MakeClients(&nodeConfig)
	if err != nil {
		fmt.Printf("Failed to make clients: %s", err)
	}
	return err
}

func rootPreFlight(ccmd *cobra.Command, args []string) error {
	if showVersion {
		fmt.Printf("algorand %s (%s)\n", version, commit)
		os.Exit(0)
	}

	ccmd.HelpFunc()(ccmd, args)
	return fmt.Errorf("")
}

func startAlgorand(ccmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	// Read configuration
	cobra.OnInitialize(readConfig)

	// Local flags
	AlgorandCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Display the application version")

	// Persistent flags
	AlgorandCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "Path to configuration file")
	AlgorandCmd.PersistentFlags().String("log-level", "INFO", "Output level of logs (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
	AlgorandCmd.PersistentFlags().StringP("host", "H", "127.0.0.1", "Algorand node hostname/IP")
	AlgorandCmd.PersistentFlags().String("algod-port", "8080", "Port used by `algod'")
	AlgorandCmd.PersistentFlags().String("algod-token", "", "Authorization token for `algod'")
	AlgorandCmd.PersistentFlags().String("kmd-port", "7833", "Port used by `kmd'")
	AlgorandCmd.PersistentFlags().String("kmd-token", "", "Authorization token for `kmd'")
	viper.BindPFlag("log-level", AlgorandCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("host", AlgorandCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("algod-port", AlgorandCmd.PersistentFlags().Lookup("algod-port"))
	viper.BindPFlag("algod-token", AlgorandCmd.PersistentFlags().Lookup("algod-token"))
	viper.BindPFlag("kmd-port", AlgorandCmd.PersistentFlags().Lookup("kmd-port"))
	viper.BindPFlag("kmd-token", AlgorandCmd.PersistentFlags().Lookup("kmd-token"))

	// Commands
	AlgorandCmd.AddCommand(statusCmd)
	AlgorandCmd.AddCommand(blockCmd)
	AlgorandCmd.AddCommand(createCmd)
	AlgorandCmd.AddCommand(backupCmd)
	AlgorandCmd.AddCommand(restoreCmd)
	AlgorandCmd.AddCommand(getCmd)
	AlgorandCmd.AddCommand(getAllCmd)
	AlgorandCmd.AddCommand(findCmd)
	AlgorandCmd.AddCommand(signCmd)

	// Hidden or aliased commands
	AlgorandCmd.AddCommand(healthCmd)
}
