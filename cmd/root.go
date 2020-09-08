package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/findy-network/findy-agent-cli/utils"
	"github.com/findy-network/findy-agent/cmds/agency"
	"github.com/lainio/err2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: utils.Version,
	Use:     "findy-agent-cli",
	Short:   "Findy agent cli tool",
	Long: `
Findy agent cli tool
	`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		agency.ParseLoggingArgs(rootFlags.logging)
	},
}

// Execute root
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// To fix errors printed twice removing the cobra generators next
		// see: https://github.com/spf13/cobra/issues/304
		//fmt.Println(err)

		os.Exit(1)
	}
}

// RootCmd returns a current root command which can be used for adding own
// commands in an own repo.
//  	implCmd.AddCommand(listCmd)
// That's a helper function to extend this CLI with own commands and offering
// same base commands as this CLI.
func RootCmd() *cobra.Command {
	return rootCmd
}

// DryRun returns a value of a dry run flag. That's a helper function to extend
// this CLI with own commands and offering same base commands as this CLI.
func DryRun() bool {
	return rootFlags.dryRun
}

// RootFlags are the common flags
type RootFlags struct {
	cfgFile string
	dryRun  bool
	logging string
}

// ClientFlags agent flags
type ClientFlags struct {
	WalletName string
	WalletKey  string
	URL        string
}

var rootFlags = RootFlags{}
var cFlags = ClientFlags{}

func init() {
	defer err2.Catch(func(err error) {
		log.Println(err)
	})

	cobra.OnInitialize(initConfig)

	flags := rootCmd.PersistentFlags()
	flags.StringVar(&rootFlags.cfgFile, "config", "", "configuration file")
	flags.StringVar(&rootFlags.logging, "logging", "-logtostderr=true -v=2", "logging startup arguments")
	flags.BoolVarP(&rootFlags.dryRun, "dry-run", "n", false, "perform a trial run with no changes made")

	err2.Check(viper.BindPFlag("logging", flags.Lookup("logging")))
	err2.Check(viper.BindPFlag("dry-run", flags.Lookup("dry-run")))
}

func initConfig() {
	viper.SetEnvPrefix("FINDY_AGENT_CLI")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv() // read in environment variables that match
	if rootFlags.cfgFile != "" {
		viper.SetConfigFile(rootFlags.cfgFile)
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
	handleViperFlags(rootCmd.Commands())
	readBoundRootFlags()
	aCmd.PreRun()
}

func readBoundRootFlags() {
	rootFlags.logging = viper.GetString("logging")
	rootFlags.dryRun = viper.GetBool("dry-run")
}

func handleViperFlags(commands []*cobra.Command) {
	for _, cmd := range commands {
		setRequiredStringFlags(cmd)
		if cmd.HasSubCommands() {
			handleViperFlags(cmd.Commands())
		}
	}
}

//TODO: change to handle all flag types
func setRequiredStringFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}

// SubCmdNeeded prints the help and error messages because the cmd is abstract.
func SubCmdNeeded(cmd *cobra.Command) {
	fmt.Println("Subcommand needed!")
	cmd.Help()
	os.Exit(1)
}
