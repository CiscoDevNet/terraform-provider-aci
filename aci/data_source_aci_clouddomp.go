package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudDomainProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciCloudDomainProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"site_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciCloudDomainProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("clouddomp")

	dn := fmt.Sprintf("uni/%s", rn)

	cloudDomP, err := getRemoteCloudDomainProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setCloudDomainProfileAttributes(cloudDomP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
