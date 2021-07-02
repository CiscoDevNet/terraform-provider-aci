package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudCIDRPool() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciCloudCIDRPoolRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_context_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"primary": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciCloudCIDRPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	addr := d.Get("addr").(string)

	rn := fmt.Sprintf("cidr-[%s]", addr)
	CloudContextProfileDn := d.Get("cloud_context_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudContextProfileDn, rn)

	cloudCidr, err := getRemoteCloudCIDRPool(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setCloudCIDRPoolAttributes(cloudCidr, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
