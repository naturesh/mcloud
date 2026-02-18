package core

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed setup.sh
var setupTemplate string

type SetupParams struct {
	VolumeName    string
	DevicePath    string
	ServerType    string
	ServerVersion string
	RconPassword  string
	DockerImage   string
}

func GenerateSetupScript(params SetupParams) (string, error) {
	tmpl, err := template.New("setup").Parse(setupTemplate)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidFormat, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", fmt.Errorf("%w: %v", ErrSetupGenerate, err)
	}

	return buf.String(), nil
}
