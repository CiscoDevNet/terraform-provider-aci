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

		Schema: AppendAttrSchemas(
			GetAnnotationAttrSchema(),
			map[string]*schema.Schema{
				"spine_profile_dn": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},

				"tdn": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
			},
		),
	}
}
