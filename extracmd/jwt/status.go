package jwt

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-cli/cmd"
	"github.com/findy-network/findy-grpc/agency/client"
	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/lainio/err2"
	"github.com/spf13/cobra"
)

var statusDoc = `    CONNECT = 1;
    ISSUE = 2;
    PROPOSE_PROOFING = 3;
    TRUST_PING = 4;
    BASIC_MESSAGE = 5;
`

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "status command for JWT gRPC",
	Long:  statusDoc,
	PreRunE: func(c *cobra.Command, args []string) (err error) {
		return cmd.BindEnvs(envs, "")
	},
	RunE: func(c *cobra.Command, args []string) (err error) {
		defer err2.Return(&err)

		if cmd.DryRun() {
			fmt.Println(dto.ToJSON(CmdData))
			return nil
		}
		c.SilenceUsage = true

		baseCfg := client.BuildClientConnBase("", CmdData.APIService, CmdData.Port, nil)
		conn = client.TryOpen(CmdData.CaDID, baseCfg)
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		didComm := agency.NewDIDCommClient(conn)
		statusResult, err := didComm.Status(ctx, &agency.ProtocolID{
			TypeId: agency.Protocol_Type(MyTypeID), // casting!!!
			Role:   agency.Protocol_RESUME,
			Id:     MyProtocolID,
		})
		err2.Check(err)

		fmt.Println("result:", statusResult.StatusJson, statusResult.State.ProtocolId.TypeId)
		return nil
	},
}

var MyTypeID int32

func init() {
	defer err2.Catch(func(err error) {
		fmt.Println(err)
	})

	statusCmd.Flags().StringVarP(&MyProtocolID, "id", "i", "", "protocol id for continue")
	statusCmd.Flags().Int32VarP(&MyTypeID, "type", "t", 3, "3=trust ping, 1=issue, ... see usage")

	JwtCmd.AddCommand(statusCmd)
}