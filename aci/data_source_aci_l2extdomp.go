package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL2Domain() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciL2DomainRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
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

func dataSourceAciL2DomainRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("l2dom-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	l2extDomP, err := getRemoteL2Domain(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setL2DomainAttributes(l2extDomP, d)
	return nil
}
