package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciEndpointSecurityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciEndpointSecurityGroupRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"application_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flood_on_encap": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"match_t": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pc_enf_pref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pref_gr_memb": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"prio": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciEndpointSecurityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	ApplicationProfileDn := d.Get("application_profile_dn").(string)
	rn := fmt.Sprintf("esg-%s", name)
	dn := fmt.Sprintf("%s/%s", ApplicationProfileDn, rn)
	fvESg, err := getRemoteEndpointSecurityGroup(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setEndpointSecurityGroupAttributes(fvESg, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
