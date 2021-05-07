package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAccessPortSelector() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciAccessPortSelectorRead,

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
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciAccessPortSelectorRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	access_port_selector_type := d.Get("access_port_selector_type").(string)

	rn := fmt.Sprintf("hports-%s-typ-%s", name, access_port_selector_type)
	LeafInterfaceProfileDn := d.Get("leaf_interface_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LeafInterfaceProfileDn, rn)

	infraHPortS, err := getRemoteAccessPortSelector(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setAccessPortSelectorAttributes(infraHPortS, d)
	return nil
}
