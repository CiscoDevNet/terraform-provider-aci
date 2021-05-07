package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3outBgpExternalPolicy() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciL3outBgpExternalPolicyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3_outside_dn": &schema.Schema{
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
		}),
	}
}

func dataSourceAciL3outBgpExternalPolicyRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("bgpExtP")
	L3OutsideDn := d.Get("l3_outside_dn").(string)

	dn := fmt.Sprintf("%s/%s", L3OutsideDn, rn)

	bgpExtP, err := getRemoteL3outBgpExternalPolicy(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setL3outBgpExternalPolicyAttributes(bgpExtP, d)
	return nil
}
