package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciFEXProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciFEXProfileRead,

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

func dataSourceAciFEXProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/fexprof-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraFexP, err := getRemoteFEXProfile(aciClient, dn)

	if err != nil {
		return err
	}
	setFEXProfileAttributes(infraFexP, d)
	return nil
}
