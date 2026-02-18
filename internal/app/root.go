package app

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ROOT = &cobra.Command{
	Use:          "mcloud",
	Short:        "mcloud, server manager",
	SilenceUsage: true,
}

func Execute(ctx context.Context) {
	if err := ROOT.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("MCLOUD")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
