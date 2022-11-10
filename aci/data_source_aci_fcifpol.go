package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciInterfaceFCPolicy() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciInterfaceFCPolicyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"automaxspeed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fill_pattern": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"port_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rx_bb_credit": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"speed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"trunk_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciInterfaceFCPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/fcIfPol-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	fcIfPol, err := getRemoteInterfaceFCPolicy(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setInterfaceFCPolicyAttributes(fcIfPol, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
