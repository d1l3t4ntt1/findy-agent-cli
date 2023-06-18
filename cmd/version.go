package cmd

import (
	"fmt"

	"github.com/findy-network/findy-agent-cli/utils"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/spf13/cobra"
)

var versionDoc = ``

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version and build information of the CLI tool",
	Long:  versionDoc,
	RunE: func(c *cobra.Command, args []string) (err error) {
		defer err2.Handle(&err)

		try.To1(fmt.Println(utils.Version))
		return nil
	},
}

func init() {
	defer err2.Catch(func(err error) {
		fmt.Println(err)
	})

	rootCmd.AddCommand(versionCmd)
}
