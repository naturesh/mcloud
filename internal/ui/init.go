package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/naturesh/mcloud/internal/core"
)

func INIT() (*core.CloudConfig, error) {
	var config core.CloudConfig

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("mcloud init").Description("Welcome! Let's setup your cloud server"),
			huh.NewInput().
				Title("Server Name").
				Placeholder("e.g. mcloud-server").
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("name must be at least 3 characters")
					}
					return nil
				}).
				Value(&config.Name),
			huh.NewSelect[string]().
				Title("Cloud Provider").
				Options(
					huh.NewOption("Digital Ocean", "digitalocean"),
				).
				Value(&config.Provider),
		),
	).WithTheme(huh.ThemeBase()).Run()
	if err != nil {
		return nil, err
	}

	options := providerOptions[config.Provider]

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Region").
				Options(options.Regions...).
				Value(&config.Region),
			huh.NewSelect[string]().
				Title("Instance Size").
				Options(options.Sizes...).
				Value(&config.InstanceSize),
			huh.NewInput().
				Title("Volume Size (GB)").
				Placeholder("e.g. 10").
				Value(&config.VolumeSize),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Server Type").
				Options(
					huh.NewOption("Paper", "PAPER"),
				).
				Value(&config.ServerOptions.Type),
			huh.NewInput().
				Title("Server Version").
				Placeholder("e.g. 1.21.11").
				Value(&config.ServerOptions.Version),
		),
	).WithTheme(huh.ThemeBase()).Run()

	if err != nil {
		return nil, err
	}

	return &config, nil
}
