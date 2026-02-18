package app

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh/spinner"
	"github.com/naturesh/mcloud/internal/core"
	"github.com/naturesh/mcloud/internal/service"
	"github.com/spf13/cobra"
)

var DOWN = &cobra.Command{
	Use:   "down [config file]",
	Short: "Stop and destroy the server instance (Data is safe)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		serv, err := service.New(args[0])
		if err != nil {
			return err
		}

		var instance core.Instance

		err = spinner.New().
			Context(ctx).
			Title(fmt.Sprintf("Searching for running server [%s]...", serv.Config.Name)).
			Type(spinner.Dots).
			ActionWithErr(func(ctx context.Context) error {
				instance, err = serv.Cloud.GetInstance(ctx, core.GetInstanceRequest{
					InstanceName: serv.Config.Name,
				})
				return err
			}).Run()

		if err != nil {
			return err
		}

		err = spinner.New().
			Context(ctx).
			Title("Saving world data safely...").
			Type(spinner.Dots).
			ActionWithErr(func(ctx context.Context) error {
				return serv.SaveData(instance.InstanceIP)
			}).Run()

		if err != nil {
			return err
		}

		err = spinner.New().
			Context(ctx).
			Title("Destroying server instance...").
			Type(spinner.Dots).
			ActionWithErr(func(ctx context.Context) error {
				return serv.Cloud.DeleteInstance(ctx, instance.InstanceID)
			}).Run()

		if err != nil {
			return err
		}

		fmt.Printf("Server stopped successfully.\n")
		fmt.Printf("- instance: Destroyed\n")
		fmt.Printf("- volume:   Preserved\n")

		return nil
	},
}

func init() {
	ROOT.AddCommand(DOWN)
}
