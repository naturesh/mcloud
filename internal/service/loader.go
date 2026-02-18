package service

import (
	"github.com/naturesh/mcloud/internal/core"
	"github.com/spf13/viper"
)

func LoadCloudConfig(path string) (core.CloudConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return core.CloudConfig{}, core.ErrConfigNotFound
	}

	var config core.CloudConfig
	if err := v.Unmarshal(&config); err != nil {
		return core.CloudConfig{}, core.ErrConfigMissing
	}

	return config, nil
}
