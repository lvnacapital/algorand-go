package cmd

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config      string // config file location
	showVersion bool   // whether to print version info or not

	// body     io.ReadWriter        // what to read/write requests body
	// verbose  bool                 // whether to display request info
	// insecure bool          = true // whether to ignore cert or not

	// key  string // blob key
	// data string // blob raw data (or '-' for stdin)
	// file string // blob file location

	// config   string // config file location
	// showVers bool   // whether to print version info or not

	// Variables to set by the linker
	version string
	commit  string

	// AlgorandCmd ...
	AlgorandCmd = &cobra.Command{
		Use:           "algorand",
		Short:         "algorand - node explorer",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		PersistentPreRunE: prePreFlight,
		PreRunE:           preFlight,
		RunE:              startAlgorand,
	}
)

func readConfig() { //(ccmd *cobra.Command, args []string) error {
	if config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name "config.yml".
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func prePreFlight(ccmd *cobra.Command, args []string) error {
	// Convert the log level.
	logLvl := lumber.LvlInt(viper.GetString("log-level"))

	// Configure the logger.
	lumber.Prefix("[algorand]")
	lumber.Level(logLvl)

	return nil
}

func preFlight(ccmd *cobra.Command, args []string) error {
	if showVersion {
		fmt.Printf("algorand %s (%s)\n", version, commit)
		os.Exit(0)
	}

	ccmd.HelpFunc()(ccmd, args)
	return fmt.Errorf("") // no error, just exit
}

func startAlgorand(ccmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	cobra.OnInitialize(readConfig)

	// Local flags
	AlgorandCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Display the application version")

	// Persistent flags
	AlgorandCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "Path to configuration file")
	AlgorandCmd.PersistentFlags().String("log-level", "INFO", "Output level of logs (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
	AlgorandCmd.PersistentFlags().StringP("host", "H", "127.0.0.1", "Algorand node hostname/IP")
	AlgorandCmd.PersistentFlags().StringP("algod-port", "a", "8080", "Port used by `algod'")
	AlgorandCmd.PersistentFlags().StringP("algod-token", "t", "", "Authorization token for `algod'")
	AlgorandCmd.PersistentFlags().StringP("kmd-port", "k", "7833", "Port used by `kmd'")
	AlgorandCmd.PersistentFlags().StringP("kmd-token", "n", "", "Authorization token for `kmd'")
	viper.BindPFlag("log-level", AlgorandCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("host", AlgorandCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("algod-port", AlgorandCmd.PersistentFlags().Lookup("algod-port"))
	viper.BindPFlag("algod-token", AlgorandCmd.PersistentFlags().Lookup("algod-token"))
	viper.BindPFlag("kmd-port", AlgorandCmd.PersistentFlags().Lookup("kmd-port"))
	viper.BindPFlag("kmd-token", AlgorandCmd.PersistentFlags().Lookup("kmd-token"))

	// Commands
	AlgorandCmd.AddCommand(statusCmd)

	// Hidden or aliased commands
	AlgorandCmd.AddCommand(healthCmd)
}
