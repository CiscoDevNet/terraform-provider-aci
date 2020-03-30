package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciVLANPool() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciVLANPoolRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"alloc_mode": &schema.Schema{
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

func dataSourceAciVLANPoolRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	allocmode := d.Get("alloc_mode").(string)

	rn := fmt.Sprintf("infra/vlanns-[%s]-%s", name, allocmode)

	dn := fmt.Sprintf("uni/%s", rn)

	fvnsVlanInstP, err := getRemoteVLANPool(aciClient, dn)

	if err != nil {
		return err
	}
	setVLANPoolAttributes(fvnsVlanInstP, d)
	return nil
}
