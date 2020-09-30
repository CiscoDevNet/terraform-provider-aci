package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciSecurityDomain() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciSecurityDomainRead,

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

func dataSourceAciSecurityDomainRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("userext/domain-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	aaaDomain, err := getRemoteSecurityDomain(aciClient, dn)

	if err != nil {
		return err
	}
	setSecurityDomainAttributes(aaaDomain, d)
	return nil
}
