package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVMMController() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciVMMControllerRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dvs_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_or_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"inventory_trig_st": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"msft_config_err_msg": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"msft_config_issues": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"n1kv_stats_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"root_cont_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"seq_num": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stats_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vxlan_depl_pref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciVMMControllerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)
	rn := fmt.Sprintf("ctrlr-%s", name)
	dn := fmt.Sprintf("%s/%s", VMMDomainDn, rn)
	vmmCtrlrP, err := getRemoteVMMController(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setVMMControllerAttributes(vmmCtrlrP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
