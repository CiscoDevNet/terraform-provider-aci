package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciAccessPortSelector() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciAccessPortSelectorRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"leaf_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"access_port_selector_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ALL",
					"range",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciAccessPortSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	access_port_selector_type := d.Get("access_port_selector_type").(string)

	rn := fmt.Sprintf("hports-%s-typ-%s", name, access_port_selector_type)
	LeafInterfaceProfileDn := d.Get("leaf_interface_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LeafInterfaceProfileDn, rn)

	infraHPortS, err := getRemoteAccessPortSelector(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	setAccessPortSelectorAttributes(infraHPortS, d)
	return nil
}
