package app

import (
	"github.com/naturesh/mcloud/internal/service"
	"github.com/spf13/cobra"
)

var STATUS = &cobra.Command{
	Use:   "status [config file]",
	Short: "Show server status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		serv, err := service.New(args[0])
		if err != nil {
			return err
		}

		err = serv.GetStatus(ctx)
		if err != nil {
			return err
		}

		return nil

	},
}

func init() {
	ROOT.AddCommand(STATUS)
}
