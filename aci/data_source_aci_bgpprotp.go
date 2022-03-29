package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3outBGPProtocolProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL3outBGPProtocolProfileRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"logical_node_profile_dn": &schema.Schema{
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

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAciL3outBGPProtocolProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	rn := fmt.Sprintf("protp")
	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)
	dn := fmt.Sprintf("%s/%s", LogicalNodeProfileDn, rn)

	bgpProtP, err := getRemoteL3outBGPProtocolProfile(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)
	_, err = setL3outBGPProtocolProfileAttributes(bgpProtP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
