package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciSpineInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciSpineInterfaceProfileRead,

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

func dataSourceAciSpineInterfaceProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/spaccportprof-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraSpAccPortP, err := getRemoteSpineInterfaceProfile(aciClient, dn)

	if err != nil {
		return err
	}
	setSpineInterfaceProfileAttributes(infraSpAccPortP, d)
	return nil
}
