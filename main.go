package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aci.Provider})
}
