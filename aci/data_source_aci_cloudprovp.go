package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudProviderProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciCloudProviderProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"vendor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func dataSourceAciCloudProviderProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	vendor := d.Get("vendor").(string)

	rn := fmt.Sprintf("clouddomp/provp-%s", vendor)

	dn := fmt.Sprintf("uni/%s", rn)

	cloudProvP, err := getRemoteCloudProviderProfile(aciClient, dn)

	if err != nil {
		return err
	}
	setCloudProviderProfileAttributes(cloudProvP, d)
	return nil
}
