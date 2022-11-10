package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3outHSRPInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL3outHSRPInterfaceProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciL3outHSRPInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("hsrpIfP")
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalInterfaceProfileDn, rn)

	hsrpIfP, err := getRemoteL3outHSRPInterfaceProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setL3outHSRPInterfaceProfileAttributes(hsrpIfP, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
