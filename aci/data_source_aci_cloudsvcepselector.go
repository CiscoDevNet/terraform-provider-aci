package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudServiceEndpointSelector() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciCloudServiceEndpointSelectorRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"cloud_service_epg_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"match_expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciCloudServiceEndpointSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	CloudServiceEPgDn := d.Get("cloud_service_epg_dn").(string)
	rn := fmt.Sprintf(models.RnCloudSvcEPSelector, name)
	dn := fmt.Sprintf("%s/%s", CloudServiceEPgDn, rn)

	cloudSvcEPSelector, err := getRemoteCloudServiceEndpointSelector(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setCloudServiceEndpointSelectorAttributes(cloudSvcEPSelector, d)
	if err != nil {
		return nil
	}

	return nil
}
