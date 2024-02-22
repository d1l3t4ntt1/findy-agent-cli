package agent

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-cli/cmd"
	"github.com/findy-network/findy-common-go/agency/client"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/spf13/cobra"
)

var invitationDoc = `Commands the cloud agent to produce invitation JSON.

If conn-id is given our end of the connection used that for naming the pairwise.
if conn-id is empty CA will genereta new UUID which will be used for both ends
of the pairwise.`

var invitationCmd = &cobra.Command{
	Use:   "invitation",
	Short: "Print connection invitation",
	Long:  invitationDoc,
	PreRunE: func(*cobra.Command, []string) (err error) {
		return cmd.BindEnvs(envs, "")
	},
	RunE: func(*cobra.Command, []string) (err error) {
		defer err2.Handle(&err)

		if cmd.DryRun() {
			fmt.Println("JWT:", CmdData.JWT)
			fmt.Println("Server:", cmd.ServiceAddr())
			fmt.Println("Label:", ourLabel)
			fmt.Println("ConnectionID:", connID)
			if connID == "" {
				fmt.Println("autogenerated shared connection ID")
			}
			return nil
		}
		baseCfg := try.To1(cmd.BaseCfg())
		conn := client.TryAuthOpen(CmdData.JWT, baseCfg)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		agent := agency.NewAgentServiceClient(conn)
		r := try.To1(agent.CreateInvitation(ctx, &agency.InvitationBase{
			ID:    connID,
			Label: ourLabel,
		}))

		if urlFormat {
			fmt.Print(r.URL)
		} else {
			fmt.Println(r.JSON)
		}

		return nil
	},
}

var (
	connID    string
	urlFormat bool
)

func init() {
	defer err2.Catch(err2.Err(func(err error) {
		fmt.Println(err)
	}))

	invitationCmd.Flags().BoolVarP(&urlFormat,
		"url", "u", false, "if set returns URL formatted invitation")
	invitationCmd.Flags().StringVar(&ourLabel,
		"label", "", "our Aries connection Label ")
	invitationCmd.Flags().StringVarP(&connID, "conn-id", "c", "",
		"connection id (UUID) for our end, if empty autogenerated for both")

	AgentCmd.AddCommand(invitationCmd)
}
