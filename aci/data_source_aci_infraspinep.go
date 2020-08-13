package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciSpineProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciSpineProfileRead,

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

func dataSourceAciSpineProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/spprof-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraSpineP, err := getRemoteSpineProfile(aciClient, dn)

	if err != nil {
		return err
	}
	setSpineProfileAttributes(infraSpineP, d)
	return nil
}
