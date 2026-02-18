package app

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh/spinner"
	"github.com/naturesh/mcloud/internal/service"
	"github.com/spf13/cobra"
)

var UP = &cobra.Command{
	Use:   "up [config file]",
	Short: "Create and start server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		serv, err := service.New(args[0])
		if err != nil {
			return err
		}

		var ip string

		err = spinner.New().
			Context(ctx).
			Title(fmt.Sprintf("Preparing server [%s]", serv.Config.Name)).
			Type(spinner.Dots).
			ActionWithErr(func(ctx context.Context) error {
				ip, err = serv.ProvisionInfrastructure(ctx)
				return err
			}).Run()

		if err != nil {
			return err
		}

		fmt.Printf("Server instance created!\n")
		fmt.Printf("- configuration: %s\n", args[0])
		fmt.Printf("- region:        %s\n", serv.Config.Region)
		fmt.Printf("- version:       %s\n", serv.Config.ServerOptions.Version)
		fmt.Printf("- type:          %s\n", serv.Config.ServerOptions.Type)
		fmt.Printf("- ip:            %s:25565\n", ip)
		fmt.Printf("\n")

		err = spinner.New().
			Context(ctx).
			Title("Waiting for os, docker, server to be ready...").
			Type(spinner.Dots).
			ActionWithErr(func(ctx context.Context) error {
				return serv.WaitForServerOpen(ip)
			}).Run()

		if err != nil {
			return err
		}

		fmt.Println("The server is open. It's accessible!")

		return nil
	},
}

func init() {
	ROOT.AddCommand(UP)
}
