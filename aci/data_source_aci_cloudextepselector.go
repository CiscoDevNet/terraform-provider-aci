package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudEndpointSelectorforExternalEPgs() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciCloudEndpointSelectorforExternalEPgsRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_external_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"is_shared": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"match_expression": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciCloudEndpointSelectorforExternalEPgsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("extepselector-%s", name)
	CloudExternalEPgDn := d.Get("cloud_external_epg_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudExternalEPgDn, rn)

	cloudExtEPSelector, err := getRemoteCloudEndpointSelectorforExternalEPgs(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	setCloudEndpointSelectorforExternalEPgsAttributes(cloudExtEPSelector, d)
	return nil
}
