package jwt

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-cli/cmd"
	"github.com/findy-network/findy-grpc/agency/client"
	"github.com/lainio/err2"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect command for JWT gRPC",
	Long: `
`,
	PreRunE: func(c *cobra.Command, args []string) (err error) {
		return cmd.BindEnvs(envs, "")
	},
	RunE: func(c *cobra.Command, args []string) (err error) {
		defer err2.Return(&err)

		if cmd.DryRun() {
			return nil
		}
		c.SilenceUsage = true

		baseCfg := client.BuildClientConnBase("", CmdData.APIService, CmdData.Port, nil)
		conn = client.TryOpen(CmdData.CaDID, baseCfg)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		pw := client.Pairwise{
			Conn:  conn,
			Label: ourLabel,
		}
		connID, ch, err := pw.Connection(ctx, invitationJSON)
		err2.Check(err)
		for status := range ch {
			fmt.Printf("Connection status: %s|%s: %s\n", connID, status.ProtocolId, status.State)
		}
		return nil
	},
}

var (
	invitationJSON string
	ourLabel       string
)

func init() {
	defer err2.Catch(func(err error) {
		fmt.Println(err)
	})

	connectCmd.Flags().StringVar(&invitationJSON, "invitation", "", "invitation json")
	connectCmd.Flags().StringVar(&ourLabel, "our-label", "", "our Aries connection Label ")

	JwtCmd.AddCommand(connectCmd)
}