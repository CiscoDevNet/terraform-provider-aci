package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciAccessGroup() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciAccessAccessGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"access_port_selector_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"fex_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciAccessAccessGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("rsaccBaseGrp")
	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	dn := fmt.Sprintf("%s/%s", AccessPortSelectorDn, rn)

	infraRsAccBaseGrp, err := getRemoteAccessAccessGroup(aciClient, dn)

	if err != nil {
		return err
	}
	setAccessAccessGroupAttributes(infraRsAccBaseGrp, d)
	return nil
}
