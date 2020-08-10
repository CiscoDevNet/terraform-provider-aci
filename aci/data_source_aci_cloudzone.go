package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciCloudAvailabilityZone() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciCloudAvailabilityZoneRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_providers_region_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

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

func dataSourceAciCloudAvailabilityZoneRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("zone-%s", name)
	CloudProvidersRegionDn := d.Get("cloud_providers_region_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudProvidersRegionDn, rn)

	cloudZone, err := getRemoteCloudAvailabilityZone(aciClient, dn)

	if err != nil {
		return err
	}
	setCloudAvailabilityZoneAttributes(cloudZone, d)
	return nil
}
