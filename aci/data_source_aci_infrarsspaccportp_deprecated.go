package aci

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciInterfaceProfileDeprecated() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use aci_spine_interface_profile_selector data source instead",
		ReadContext:        dataSourceAciInterfaceProfileRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"spine_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
