package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciBGPTimersPolicy() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciBGPTimersPolicyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"gr_ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hold_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ka_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_as_limit": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"stale_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciBGPTimersPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("bgpCtxP-%s", name)
	TenantDn := d.Get("tenant_dn").(string)

	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	bgpCtxPol, err := getRemoteBGPTimersPolicy(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)
	_, err = setBGPTimersPolicyAttributes(bgpCtxPol, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
