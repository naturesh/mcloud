package service

import (
	"context"
	"fmt"
	"time"

	"github.com/naturesh/mcloud/internal/core"
	"github.com/naturesh/mcloud/internal/provider"
	"github.com/naturesh/mcloud/internal/ssh"
)

type Service struct {
	Config core.CloudConfig
	Cloud  core.CloudProvider
}

func New(path string) (*Service, error) {
	config, err := LoadCloudConfig(path)
	if err != nil {
		return nil, err
	}

	cloud, err := provider.New(config.Provider)
	if err != nil {
		return nil, err
	}

	return &Service{
		Config: config,
		Cloud:  cloud,
	}, nil
}

func (s *Service) ProvisionInfrastructure(ctx context.Context) (string, error) {
	var volume core.Volume
	var instance core.Instance

	volume, err := s.Cloud.GetVolume(ctx, core.GetVolumeRequest{
		VolumeName: s.Config.Name,
	})
	if err != nil {
		return "", err
	}

	if volume.VolumeID == "" {
		volume, err = s.Cloud.CreateVolume(ctx, core.CreateVolumeRequest{
			VolumeName:    s.Config.Name,
			VolumeSize:    s.Config.VolumeSize,
			Region:        s.Config.Region,
			ServerOptions: s.Config.ServerOptions,
		})
		if err != nil {
			return "", err
		}
	}

	instance, err = s.Cloud.GetInstance(ctx, core.GetInstanceRequest{
		InstanceName: s.Config.Name,
	})
	if err != nil {
		return "", err
	}

	if instance.InstanceIP != "" {
		return instance.InstanceIP, nil
	}

	instance, err = s.Cloud.CreateInstance(ctx, core.CreateInstanceRequest{
		VolumeID:      volume.VolumeID,
		VolumeName:    s.Config.Name,
		Region:        s.Config.Region,
		InstanceSize:  s.Config.InstanceSize,
		ServerOptions: s.Config.ServerOptions,
	})
	if err != nil {
		return "", err
	}

	ip, err := s.Cloud.Wait(ctx, instance.InstanceID)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func (s *Service) WaitForServerOpen(ip string) error {
	time.Sleep(20 * time.Second)

	client, err := ssh.Connect(ip)
	if err != nil {
		return err
	}

	defer client.Close()

	err = client.WaitForLog("Done", 300)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SaveData(instance_ip string) error {
	client, err := ssh.Connect(instance_ip)
	if err != nil {
		return err
	}

	defer client.Close()

	if !client.IsContainerRunning() {
		return nil
	}

	err = client.Run("docker exec mcloud rcon-cli save-all > /dev/null 2>&1")
	if err != nil {
		return err
	}

	err = client.WaitForLog("Saved the game", 60)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	err = client.Run("docker stop -t 30 mcloud > /dev/null 2>&1")
	if err != nil {
		return err
	}

	err = client.Run("sync")
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	return nil
}

func (s *Service) SendCommand(ctx context.Context, command string) error {
	instance, err := s.Cloud.GetInstance(ctx, core.GetInstanceRequest{
		InstanceName: s.Config.Name,
	})
	if err != nil {
		return err
	}

	if instance.InstanceID == "" {
		return fmt.Errorf("%w: server:%s", core.ErrNotFound, s.Config.Name)
	}

	client, err := ssh.Connect(instance.InstanceIP)
	if err != nil {
		return err
	}

	defer client.Close()

	return client.Run(fmt.Sprintf("docker exec mcloud rcon-cli '%s'", command))
}

func (s *Service) GetStatus(ctx context.Context) error {
	instance, err := s.Cloud.GetInstance(ctx, core.GetInstanceRequest{
		InstanceName: s.Config.Name,
	})
	if err != nil {
		return err
	}
	if instance.InstanceID == "" {
		return fmt.Errorf("%w: server:%s", core.ErrNotFound, s.Config.Name)
	}

	client, err := ssh.Connect(instance.InstanceIP)
	if err != nil {
		return err
	}

	defer client.Close()

	cmd := fmt.Sprintf(`
		echo "- ip:       %s"

		INSTANCE_STATS=$(docker stats --no-stream --format "{{.CPUPerc}} {{.MemUsage}}" mcloud 2>/dev/null)
		VOLUME_STATS=$(df -h /mnt/%s | awk 'NR==2 {print $5, $4}')

		if [ -z "$INSTANCE_STATS" ]; then
			echo "- instance: stopped"
		else
			set -- $INSTANCE_STATS
			echo "- instance: CPU $1 | RAM $2/$4"
		fi

		if [ -z "$VOLUME_STATS" ]; then
			echo "- volume:   not mounted"
		else
			set -- $VOLUME_STATS
			echo "- volume:   Used $1 | Free $2"
		fi
	`, instance.InstanceIP, s.Config.Name)
	if err = client.Run(cmd); err != nil {
		return err
	}

	return nil
}
