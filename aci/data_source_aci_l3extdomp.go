package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciL3DomainProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciL3DomainProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
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

func dataSourceAciL3DomainProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("l3dom-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	l3extDomP, err := getRemoteL3DomainProfile(aciClient, dn)

	if err != nil {
		return err
	}
	setL3DomainProfileAttributes(l3extDomP, d)
	return nil
}
