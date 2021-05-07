package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudProvidersRegion() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciCloudProvidersRegionRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_provider_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"admin_st": &schema.Schema{
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

func dataSourceAciCloudProvidersRegionRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("region-%s", name)
	CloudProviderProfileDn := d.Get("cloud_provider_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudProviderProfileDn, rn)

	cloudRegion, err := getRemoteCloudProvidersRegion(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setCloudProvidersRegionAttributes(cloudRegion, d)
	return nil
}
