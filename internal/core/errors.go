package core

import "errors"

var (
	ErrConfigMissing  = errors.New("configuration file is missing")
	ErrConfigNotFound = errors.New("configuration file not found")
	ErrSetupGenerate  = errors.New("failed to generate setup.sh")

	ErrInvalidFormat     = errors.New("invalid format")
	ErrSSHKeyList        = errors.New("failed to list ssh keys")
	ErrInstanceCreate    = errors.New("failed to create instance")
	ErrInstanceTagCreate = errors.New("failed to tag instance")
	ErrInstanceList      = errors.New("failed to list instance")
	ErrInstanceDelete    = errors.New("failed to delete instance")
	ErrVolumeCreate      = errors.New("failed to create volume")
	ErrVolumeTagCreate   = errors.New("failed to tag volume")
	ErrVolumeList        = errors.New("failed to list volume")
	ErrVolumeDelete      = errors.New("failed to delete volume")
	ErrWaitLoadIP        = errors.New("failed to load ip")
	ErrTagCreate         = errors.New("failed to create tag")
	ErrTagApply          = errors.New("failed to apply tag")

	ErrTokenNotFound       = errors.New("provider token not found")
	ErrUnsupportedProvider = errors.New("unsupported provider")

	ErrNotFound = errors.New("not found")

	ErrSSHAgentNotFound = errors.New("ssh-agent not found")
	ErrSSHConnect       = errors.New("failed to connect ssh")

	ErrMarshal   = errors.New("failed to marshal configuration file")
	ErrWriteFile = errors.New("failed to write file")
)
