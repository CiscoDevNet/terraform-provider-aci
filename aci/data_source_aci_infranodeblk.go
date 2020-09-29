package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciNodeBlock() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciNodeBlockRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"switch_association_dn": &schema.Schema{
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

func dataSourceAciNodeBlockRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("nodeblk-%s", name)
	SwitchAssociationDn := d.Get("switch_association_dn").(string)

	dn := fmt.Sprintf("%s/%s", SwitchAssociationDn, rn)

	infraNodeBlk, err := getRemoteNodeBlock(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setNodeBlockAttributes(infraNodeBlk, d)
	return nil
}
