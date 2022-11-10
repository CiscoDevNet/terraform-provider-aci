package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudProviderProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciCloudProviderProfileRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{

			"vendor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAciCloudProviderProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	vendor := d.Get("vendor").(string)

	rn := fmt.Sprintf("clouddomp/provp-%s", vendor)

	dn := fmt.Sprintf("uni/%s", rn)

	cloudProvP, err := getRemoteCloudProviderProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setCloudProviderProfileAttributes(cloudProvP, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getRemoteCloudProviderProfile(client *client.Client, dn string) (*models.CloudProviderProfile, error) {
	cloudProvPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudProvP := models.CloudProviderProfileFromContainer(cloudProvPCont)

	if cloudProvP.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudProviderProfile %s not found", cloudProvP.DistinguishedName)
	}

	return cloudProvP, nil
}

func setCloudProviderProfileAttributes(cloudProvP *models.CloudProviderProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cloudProvP.DistinguishedName)
	cloudProvPMap, err := cloudProvP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("annotation", cloudProvPMap["annotation"])
	d.Set("vendor", cloudProvPMap["vendor"])
	return d, nil
}
