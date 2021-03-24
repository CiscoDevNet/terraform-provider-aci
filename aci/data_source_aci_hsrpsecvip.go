package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciL3outHSRPSecondaryVIP() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciL3outHSRPSecondaryVIPRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"hsrp_group_profile_dn": &schema.Schema{
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

func dataSourceAciL3outHSRPSecondaryVIPRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	ip := d.Get("ip").(string)

	rn := fmt.Sprintf("hsrpSecVip-[%s]", ip)
	HSRPGroupProfileDn := d.Get("hsrp_group_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", HSRPGroupProfileDn, rn)

	hsrpSecVip, err := getRemoteL3outHSRPSecondaryVIP(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setL3outHSRPSecondaryVIPAttributes(hsrpSecVip, d)
	return nil
}
