package app

import (
	"fmt"
	"os"

	"github.com/naturesh/mcloud/internal/core"
	"github.com/naturesh/mcloud/internal/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var INIT = &cobra.Command{
	Use:   "init [config file]",
	Short: "Create a new server configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		config, err := ui.INIT()

		if err != nil {
			return err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("%w: %v", core.ErrMarshal, err)
		}

		if err := os.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("%w: %v", core.ErrWriteFile, err)
		}

		fmt.Println("mcloud init completed.")

		return nil
	},
}

func init() {
	ROOT.AddCommand(INIT)
}
