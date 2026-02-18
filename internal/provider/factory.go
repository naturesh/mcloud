package provider

import (
	"github.com/naturesh/mcloud/internal/core"
	"github.com/naturesh/mcloud/internal/provider/digitalocean"
	"github.com/spf13/viper"
)

func New(providerType string) (core.CloudProvider, error) {
	switch providerType {
	case "digitalocean":
		token := viper.GetString("digitalocean.token")
		if token == "" {
			return nil, core.ErrTokenNotFound
		}
		return digitalocean.New(token), nil
	default:
		return nil, core.ErrUnsupportedProvider
	}
}
