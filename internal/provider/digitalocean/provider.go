package digitalocean

import (
	"context"
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/naturesh/mcloud/internal/core"
)

var _ core.CloudProvider = (*DigitalOcean)(nil)

type DigitalOcean struct {
	client *godo.Client
}

func New(token string) *DigitalOcean {
	return &DigitalOcean{
		client: godo.NewFromToken(token),
	}
}

func (p *DigitalOcean) CreateInstance(ctx context.Context, req core.CreateInstanceRequest) (core.Instance, error) {
	keys, _, err := p.client.Keys.List(ctx, nil)
	if err != nil {
		return core.Instance{}, fmt.Errorf("%w: %v", core.ErrSSHKeyList, err)
	}

	var keyRefs []godo.DropletCreateSSHKey
	for _, k := range keys {
		keyRefs = append(keyRefs, godo.DropletCreateSSHKey{ID: k.ID})
	}

	setup, err := core.GenerateSetupScript(core.SetupParams{
		VolumeName:    req.VolumeName,
		DevicePath:    fmt.Sprintf("/dev/disk/by-id/scsi-0DO_Volume_%s", req.VolumeName),
		ServerType:    req.ServerOptions.Type,
		ServerVersion: req.ServerOptions.Version,
		RconPassword:  core.DefaultRconPassword,
		DockerImage:   core.DefaultDockerImage,
	})
	if err != nil {
		return core.Instance{}, err
	}

	request := &godo.DropletCreateRequest{
		Name:   req.VolumeName,
		Region: req.Region,
		Size:   req.InstanceSize,
		Image:  godo.DropletCreateImage{Slug: core.DefaultUbuntuImage},
		Volumes: []godo.DropletCreateVolume{
			{ID: req.VolumeID},
		},
		SSHKeys:  keyRefs,
		UserData: setup,
	}

	instance, _, err := p.client.Droplets.Create(ctx, request)
	if err != nil {
		return core.Instance{}, fmt.Errorf("%w: %v", core.ErrInstanceCreate, err)
	}

	err = applyTags(
		ctx,
		p,
		strconv.Itoa(instance.ID),
		godo.DropletResourceType,
		req.ServerOptions,
	)
	if err != nil {
		return core.Instance{}, fmt.Errorf("%w: %v", core.ErrInstanceTagCreate, err)
	}

	return core.Instance{
		InstanceID:    strconv.Itoa(instance.ID),
		InstanceName:  req.VolumeName,
		InstanceIP:    "",
		ServerOptions: req.ServerOptions,
	}, nil
}

func (p *DigitalOcean) ListInstances(ctx context.Context) ([]core.Instance, error) {
	droplets, _, err := p.client.Droplets.ListByTag(ctx, core.DefaultTag, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", core.ErrInstanceList, err)
	}

	var instances []core.Instance

	for _, d := range droplets {
		ip, _ := d.PublicIPv4()
		opts := parseServerOptions(d.Tags)

		instances = append(instances, core.Instance{
			InstanceID:    strconv.Itoa(d.ID),
			InstanceName:  d.Name,
			InstanceIP:    ip,
			ServerOptions: opts,
		})
	}

	return instances, nil
}

func (p *DigitalOcean) GetInstance(ctx context.Context, req core.GetInstanceRequest) (core.Instance, error) {
	instances, err := p.ListInstances(ctx)
	if err != nil {
		return core.Instance{}, err
	}

	for _, i := range instances {
		if i.InstanceName == req.InstanceName {
			return i, nil
		}
	}

	return core.Instance{}, err
}

func (p *DigitalOcean) DeleteInstance(ctx context.Context, instance_id string) error {
	i_id, err := strconv.Atoi(instance_id)
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrInvalidFormat, err)
	}

	_, err = p.client.Droplets.Delete(ctx, i_id)
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrInstanceDelete, err)
	}

	return nil
}

func (p *DigitalOcean) CreateVolume(ctx context.Context, req core.CreateVolumeRequest) (core.Volume, error) {
	i_size, err := strconv.Atoi(req.VolumeSize)
	if err != nil {
		return core.Volume{}, fmt.Errorf("%w: %v", core.ErrInvalidFormat, err)
	}

	request := &godo.VolumeCreateRequest{
		Name:           req.VolumeName,
		Region:         req.Region,
		SizeGigaBytes:  int64(i_size),
		FilesystemType: core.DefaultFileSystemType,
	}

	volume, _, err := p.client.Storage.CreateVolume(ctx, request)
	if err != nil {
		return core.Volume{}, fmt.Errorf("%w: %v", core.ErrVolumeCreate, err)
	}

	err = applyTags(
		ctx,
		p,
		volume.ID,
		godo.VolumeResourceType,
		req.ServerOptions,
	)
	if err != nil {
		return core.Volume{}, fmt.Errorf("%w: %v", core.ErrVolumeTagCreate, err)
	}

	return core.Volume{
		VolumeID:      volume.ID,
		VolumeName:    req.VolumeName,
		VolumeSize:    req.VolumeSize,
		Region:        req.Region,
		ServerOptions: req.ServerOptions,
	}, nil
}

func (p *DigitalOcean) ListVolumes(ctx context.Context) ([]core.Volume, error) {
	vols, _, err := p.client.Storage.ListVolumes(ctx, &godo.ListVolumeParams{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", core.ErrVolumeList, err)
	}

	var volumes []core.Volume

	if len(vols) > 0 {
		for _, v := range vols {
			if slices.Contains(v.Tags, core.DefaultTag) {
				volumes = append(volumes, core.Volume{
					VolumeID:      v.ID,
					VolumeName:    v.Name,
					VolumeSize:    strconv.Itoa(int(v.SizeGigaBytes)),
					Region:        v.Region.Name,
					ServerOptions: parseServerOptions(v.Tags),
				})
			}
		}
	}

	return volumes, nil
}

func (p *DigitalOcean) GetVolume(ctx context.Context, req core.GetVolumeRequest) (core.Volume, error) {
	vols, err := p.ListVolumes(ctx)
	if err != nil {
		return core.Volume{}, err
	}

	for _, v := range vols {
		if v.VolumeName == req.VolumeName {
			return v, nil
		}
	}

	return core.Volume{}, nil
}

func (p *DigitalOcean) DeleteVolume(ctx context.Context, volume_id string) error {
	_, err := p.client.Storage.DeleteVolume(ctx, volume_id)
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrVolumeDelete, err)
	}

	return nil
}

func (p *DigitalOcean) Wait(ctx context.Context, instance_id string) (string, error) {
	i_id, err := strconv.Atoi(instance_id)
	if err != nil {
		return "", fmt.Errorf("%w: %v", core.ErrInvalidFormat, err)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("%w: %v", core.ErrWaitLoadIP, ctx.Err())
		case <-ticker.C:
			droplet, _, err := p.client.Droplets.Get(ctx, i_id)
			if err != nil {
				continue
			}

			if ip, _ := droplet.PublicIPv4(); ip != "" {
				return ip, nil
			}
		}
	}
}

func createTags(option core.ServerOptions) []string {
	return []string{
		core.DefaultTag,
		fmt.Sprintf("%s%s", core.DefaultVersionTagPrefix, strings.ReplaceAll(option.Version, ".", "_")),
		fmt.Sprintf("%s%s", core.DefaultTypeTagPrefix, option.Type),
	}
}

func parseServerOptions(tags []string) core.ServerOptions {
	opts := core.ServerOptions{}

	for _, tag := range tags {
		if v, ok := strings.CutPrefix(tag, core.DefaultVersionTagPrefix); ok {
			opts.Version = strings.ReplaceAll(v, "_", ".")
		}
		if v, ok := strings.CutPrefix(tag, core.DefaultTypeTagPrefix); ok {
			opts.Type = v
		}
	}

	return opts
}

func applyTags(ctx context.Context, p *DigitalOcean, resourceID string, resourceType godo.ResourceType, option core.ServerOptions) error {
	tags := createTags(option)

	for _, tag := range tags {
		_, _, err := p.client.Tags.Create(ctx, &godo.TagCreateRequest{
			Name: tag,
		})
		if err != nil {
			return fmt.Errorf("%w: %v", core.ErrTagCreate, err)
		}

		tagRequest := &godo.TagResourcesRequest{
			Resources: []godo.Resource{{ID: resourceID, Type: resourceType}},
		}
		_, err = p.client.Tags.TagResources(ctx, tag, tagRequest)
		if err != nil {
			return fmt.Errorf("%w: %s: %v", core.ErrTagApply, resourceType, err)
		}
	}

	return nil
}
