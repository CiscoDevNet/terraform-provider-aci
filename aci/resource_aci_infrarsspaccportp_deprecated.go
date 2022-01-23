package aci

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciInterfaceProfileDeprecated() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use aci_spine_interface_profile_selector resource instead",
		CreateContext:      resourceAciInterfaceProfileCreate,
		UpdateContext:      resourceAciInterfaceProfileUpdate,
		ReadContext:        resourceAciInterfaceProfileRead,
		DeleteContext:      resourceAciInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"spine_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Default:  "orchestrator:terraform",
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
