package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudProvidersRegion() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciCloudProvidersRegionRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
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
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "orchestrator:terraform",
			},
			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAciCloudProvidersRegionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("region-%s", name)
	CloudProviderProfileDn := d.Get("cloud_provider_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudProviderProfileDn, rn)

	cloudRegion, err := getRemoteCloudProvidersRegion(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	setCloudProvidersRegionAttributes(cloudRegion, d)
	return nil
}
func getRemoteCloudProvidersRegion(client *client.Client, dn string) (*models.CloudProvidersRegion, error) {
	cloudRegionCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudRegion := models.CloudProvidersRegionFromContainer(cloudRegionCont)

	if cloudRegion.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudProvidersRegion %s not found", cloudRegion.DistinguishedName)
	}

	return cloudRegion, nil
}

func setCloudProvidersRegionAttributes(cloudRegion *models.CloudProvidersRegion, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(cloudRegion.DistinguishedName)
	//d.Set("description", cloudRegion.Description)
	if dn != cloudRegion.DistinguishedName {
		d.Set("cloud_provider_profile_dn", "")
	}
	d.Set("description", cloudRegion.Description)
	cloudRegionMap, _ := cloudRegion.ToMap()
	d.Set("annotation", cloudRegionMap["annotation"])
	d.Set("name", cloudRegionMap["name"])
	d.Set("admin_st", cloudRegionMap["adminSt"])
	d.Set("name_alias", cloudRegionMap["nameAlias"])
	return d
}
