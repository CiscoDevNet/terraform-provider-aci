package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciIGMPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciIGMPInterfacePolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"grp_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"if_ctrl": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"last_mbr_cnt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_mbr_resp_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"querier_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_intvl": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"robust_fac": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rsp_intvl": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_query_cnt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_query_intvl": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ver": {
				Type:     schema.TypeString,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciIGMPInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnIgmpIfPol, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	igmpIfPol, err := getRemoteIGMPInterfacePolicy(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setIGMPInterfacePolicyAttributes(igmpIfPol, d)
	if err != nil {
		return nil
	}

	return nil
}
