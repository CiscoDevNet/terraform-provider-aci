package main

import (
	"github.com/ciscoecosystem/terraform-provider-aci/aci"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aci.Provider})
}
