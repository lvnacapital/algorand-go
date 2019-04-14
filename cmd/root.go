package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
	binary      string

	algodClient algod.Client
	kmdClient   kmd.Client

	// Subcommand variables
	blockNumber uint64
	fromAddr    string
	toAddr      string
	noteText    string
	fee         uint64
	amount      uint64
	firstRound  uint64
	lastRound   uint64
	addr        string
	txID        string
	// WalletName ...
	WalletName string
	// WalletPassword ...
	WalletPassword string
	// WalletMnemonic ...
	WalletMnemonic string

	// Linker variables
	version string
	commit  string

	// AlgorandCmd defines the top-level command
	AlgorandCmd = &cobra.Command{
		Use:               binary,
		Short:             fmt.Sprintf("%s - Interactive Algorand node explorer.", binary),
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
		// Use config file passed in the flag
		viper.SetConfigFile(config)
	} else {
		// Find 'home' directory
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Find 'pwd'
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in 'home' or 'pwd' directory with name 'config.yml'
		viper.AddConfigPath(dir)
		viper.AddConfigPath(pwd)
		if os.Getenv("GOTEST") == "true" {
			// During testing 'os.Getwd()' might not be not root
			viper.AddConfigPath(filepath.Join(pwd, ".."))
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
	}

	if err := viper.ReadInConfig(); err != nil {
		// Falling back on environment variables
		// fmt.Println(err)
	}
}

func allPreFlight(ccmd *cobra.Command, args []string) (err error) {
	// Convert the log level.
	// logLvl := lumber.LvlInt(viper.GetString("log-level"))

	// Configure the logger.
	// lumber.Prefix("[algorand]")
	// lumber.Level(logLvl)

	nodeConfig := util.Node{
		AlgodAddress: fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("algod-port")),
		KmdAddress:   fmt.Sprintf("http://%s:%s", viper.GetString("host"), viper.GetString("kmd-port")),
		AlgodToken:   viper.GetString("algod-token"),
		KmdToken:     viper.GetString("kmd-token"),
	}

	if algodClient, kmdClient, err = util.MakeClients(&nodeConfig); err != nil {
		fmt.Printf("Failed to make clients: %s", err)
	}
	return err
}

func rootPreFlight(ccmd *cobra.Command, args []string) error {
	if showVersion {
		fmt.Printf("%s %s (%s)\n", binary, version, commit)
	} else {
		ccmd.HelpFunc()(ccmd, args)
	}

	return nil
}

func startAlgorand(ccmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	// Read configuration
	cobra.OnInitialize(readConfig)

	// Local flags
	AlgorandCmd.Flags().BoolVar(&showVersion, "version", false, "Display the application version")

	// Persistent flags
	AlgorandCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "Path to configuration FILE")
	// AlgorandCmd.PersistentFlags().String("log-level", "INFO", "Output level of logs (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
	AlgorandCmd.PersistentFlags().StringP("host", "H", "127.0.0.1", "Algorand node hostname/IP")
	AlgorandCmd.PersistentFlags().String("algod-port", "8080", "Port used by 'algod'")
	AlgorandCmd.PersistentFlags().String("algod-token", "", "Authorization token for 'algod'")
	AlgorandCmd.PersistentFlags().String("kmd-port", "7833", "Port used by 'kmd'")
	AlgorandCmd.PersistentFlags().String("kmd-token", "", "Authorization token for 'kmd'")
	// viper.BindPFlag("log-level", AlgorandCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("host", AlgorandCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("algod-port", AlgorandCmd.PersistentFlags().Lookup("algod-port"))
	viper.BindPFlag("algod-token", AlgorandCmd.PersistentFlags().Lookup("algod-token"))
	viper.BindPFlag("kmd-port", AlgorandCmd.PersistentFlags().Lookup("kmd-port"))
	viper.BindPFlag("kmd-token", AlgorandCmd.PersistentFlags().Lookup("kmd-token"))
	viper.BindEnv("host", "ALGORAND_HOST")
	viper.BindEnv("algod-port", "ALGOD_PORT")
	viper.BindEnv("algod-token", "ALGOD_TOKEN")
	viper.BindEnv("kmd-port", "KMD_PORT")
	viper.BindEnv("kmd-token", "KMD_TOKEN")

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

	// Behavior options
	AlgorandCmd.SetUsageTemplate(usageTemplate)
	AlgorandCmd.DisableSuggestions = true
}

var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} <command>{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use '{{.CommandPath}} <command> --help' for more information about a command.{{end}}
`
