package aci

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciBgpRouteControlProfile() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use aci_route_control_profile resource instead",

		CreateContext: resourceAciRouteControlProfileCreate,
		UpdateContext: resourceAciRouteControlProfileUpdate,
		ReadContext:   resourceAciRouteControlProfileRead,
		DeleteContext: resourceAciRouteControlProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRouteControlProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"route_control_profile_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"global",
					"combinable",
				}, false),
			},
		}),
	}
}
