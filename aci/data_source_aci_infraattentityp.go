package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAttachableAccessEntityProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciAttachableAccessEntityProfileRead,

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

func dataSourceAciAttachableAccessEntityProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/attentp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraAttEntityP, err := getRemoteAttachableAccessEntityProfile(aciClient, dn)

	if err != nil {
		return err
	}
	setAttachableAccessEntityProfileAttributes(infraAttEntityP, d)
	return nil
}
