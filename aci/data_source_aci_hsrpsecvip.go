package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3outHSRPSecondaryVIP() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL3outHSRPSecondaryVIPRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3out_hsrp_interface_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"config_issues": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciL3outHSRPSecondaryVIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	ip := d.Get("ip").(string)

	rn := fmt.Sprintf("hsrpSecVip-[%s]", ip)
	HSRPGroupProfileDn := d.Get("l3out_hsrp_interface_group_dn").(string)

	dn := fmt.Sprintf("%s/%s", HSRPGroupProfileDn, rn)

	hsrpSecVip, err := getRemoteL3outHSRPSecondaryVIP(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setL3outHSRPSecondaryVIPAttributes(hsrpSecVip, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
