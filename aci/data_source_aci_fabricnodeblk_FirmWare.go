package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciNodeBlockFW() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciNodeBlockReadFW,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"firmware_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"from_": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"to_": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciNodeBlockReadFW(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("nodeblk-%s", name)
	FirmwareGroupDn := d.Get("firmware_group_dn").(string)

	dn := fmt.Sprintf("%s/%s", FirmwareGroupDn, rn)

	fabricNodeBlk, err := getRemoteNodeBlockFW(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setNodeBlockAttributesFW(fabricNodeBlk, d)
	return nil
}
