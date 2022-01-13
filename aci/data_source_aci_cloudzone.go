package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudAvailabilityZone() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciCloudAvailabilityZoneRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{

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

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},
		},
	}
}

func dataSourceAciCloudAvailabilityZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("zone-%s", name)
	CloudProvidersRegionDn := d.Get("cloud_providers_region_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudProvidersRegionDn, rn)

	cloudZone, err := getRemoteCloudAvailabilityZone(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setCloudAvailabilityZoneAttributes(cloudZone, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
func getRemoteCloudAvailabilityZone(client *client.Client, dn string) (*models.CloudAvailabilityZone, error) {
	cloudZoneCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudZone := models.CloudAvailabilityZoneFromContainer(cloudZoneCont)

	if cloudZone.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudAvailabilityZone %s not found", cloudZone.DistinguishedName)
	}

	return cloudZone, nil
}

func setCloudAvailabilityZoneAttributes(cloudZone *models.CloudAvailabilityZone, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudZone.DistinguishedName)
	if dn != cloudZone.DistinguishedName {
		d.Set("cloud_providers_region_dn", "")
	}
	cloudZoneMap, err := cloudZone.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", cloudZoneMap["name"])
	d.Set("annotation", cloudZoneMap["annotation"])
	d.Set("name_alias", cloudZoneMap["nameAlias"])
	return d, nil
}
