package core

import "context"

type CloudConfig struct {
	Provider      string        `mapstructure:"provider" yaml:"provider"`
	Name          string        `mapstructure:"name" yaml:"name"`
	Region        string        `mapstructure:"region" yaml:"region"`
	InstanceSize  string        `mapstructure:"instance_size" yaml:"instance_size"`
	VolumeSize    string        `mapstructure:"volume_size" yaml:"volume_size"`
	ServerOptions ServerOptions `mapstructure:"server_options" yaml:"server_options"`
}

type ServerOptions struct {
	Type    string `mapstructure:"type" yaml:"type"`
	Version string `mapstructure:"version" yaml:"version"`
}

type Instance struct {
	InstanceID    string
	InstanceName  string
	InstanceIP    string
	ServerOptions ServerOptions
}

type Volume struct {
	VolumeID      string
	VolumeName    string
	VolumeSize    string
	Region        string
	ServerOptions ServerOptions
}

type CreateInstanceRequest struct {
	VolumeID      string
	VolumeName    string
	Region        string
	InstanceSize  string
	ServerOptions ServerOptions
}

type CreateVolumeRequest struct {
	VolumeName    string
	VolumeSize    string
	Region        string
	ServerOptions ServerOptions
}

type GetInstanceRequest struct {
	InstanceName string
}

type GetVolumeRequest struct {
	VolumeName string
}

type CloudProvider interface {
	CreateInstance(ctx context.Context, req CreateInstanceRequest) (Instance, error)
	Wait(ctx context.Context, instance_id string) (string, error)
	ListInstances(ctx context.Context) ([]Instance, error)
	GetInstance(ctx context.Context, req GetInstanceRequest) (Instance, error)
	DeleteInstance(ctx context.Context, instance_id string) error

	CreateVolume(ctx context.Context, req CreateVolumeRequest) (Volume, error)
	ListVolumes(ctx context.Context) ([]Volume, error)
	GetVolume(ctx context.Context, req GetVolumeRequest) (Volume, error)
	DeleteVolume(ctx context.Context, volume_id string) error
}
