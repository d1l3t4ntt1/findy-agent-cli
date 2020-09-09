package cmd

import (
	"log"
	"os"

	"github.com/findy-network/findy-agent/cmds/key"
	"github.com/lainio/err2"
	"github.com/spf13/cobra"
)

// keyCmd represents the key subcommand
var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Parent command for handling keys",
	Long: `
Parent command for handling keys
	`,
	Run: func(cmd *cobra.Command, args []string) {
		SubCmdNeeded(cmd)
	},
}

// createKeyCmd represents the createkey subcommand
var createKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Command for creating valid wallet keys",
	Long: `
Command for creating valid wallet keys	

Example	
	findy-agent-cli tools key create \
		--seed 00000000000000000000thisisa_test
	`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defer err2.Return(&err)
		err2.Check(keyCreateCmd.Validate())
		if !rootFlags.dryRun {
			cmd.SilenceUsage = true
			err2.Try(keyCreateCmd.Exec(os.Stdout))
		}
		return nil
	},
}

var keyCreateCmd = key.CreateCmd{}

func init() {
	defer err2.Catch(func(err error) {
		log.Println(err)
	})

	createKeyCmd.Flags().StringVar(&keyCreateCmd.Seed, "seed", "", "seed for wallet key creation")

	toolsCmd.AddCommand(keyCmd)
	keyCmd.AddCommand(createKeyCmd)
}
