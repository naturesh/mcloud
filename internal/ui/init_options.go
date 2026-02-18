package ui

import "github.com/charmbracelet/huh"


var providerOptions = map[string]struct{
	Regions []huh.Option[string]
	Sizes   []huh.Option[string]
}{
	"digitalocean" : {
		Regions: []huh.Option[string]{
			huh.NewOption("New York City", "nyc1"),
			huh.NewOption("Amsterdam", "ams3"),
			huh.NewOption("San Francisco", "sfo3"),
			huh.NewOption("Singapore", "sgp1"),
			huh.NewOption("London", "lon1"),
			huh.NewOption("Frankfurt", "fra1"),
			huh.NewOption("Toronto", "tor1"),
			huh.NewOption("Bangalore", "blr1"),
			huh.NewOption("Sydney", "syd1"),
		},
		Sizes: []huh.Option[string]{
			huh.NewOption("s-1vcpu-2gb-intel       / $14/mo", "s-1vcpu-2gb-intel"),
			huh.NewOption("s-2vcpu-2gb-intel       / $21/mo", "s-2vcpu-2gb-intel"),
			huh.NewOption("s-2vcpu-4gb-intel       / $28/mo", "s-2vcpu-4gb-intel"),
			huh.NewOption("s-2vcpu-8gb-160gb-intel / $48/mo", "s-2vcpu-8gb-160gb-intel"),
			huh.NewOption("s-4vcpu-8gb-intel       / $56/mo", "s-4vcpu-8gb-intel"),
			huh.NewOption("s-8vcpu-16gb-intel      / $112/mo", "s-8vcpu-16gb-intel"),
		},
	},
}