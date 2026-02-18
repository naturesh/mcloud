package app

import (
	"strings"

	"github.com/naturesh/mcloud/internal/service"
	"github.com/spf13/cobra"
)

var CONSOLE = &cobra.Command{
	Use:   "console [config file] [command...]",
	Short: "Send a command to the Minecraft server console",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		serv, err := service.New(args[0])
		if err != nil {
			return err
		}

		command := strings.Join(args[1:], " ")

		err = serv.SendCommand(ctx, command)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	ROOT.AddCommand(CONSOLE)
}
