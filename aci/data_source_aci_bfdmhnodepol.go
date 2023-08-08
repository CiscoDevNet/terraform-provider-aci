package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciBFDMultihopNodePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciBFDMultihopNodePolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"admin_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"detection_multiplier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_rx_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_tx_interval": {
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

func dataSourceAciBFDMultihopNodePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnBfdMhNodePol, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	bfdMhNodePol, err := getRemoteBFDMultihopNodePolicy(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setBFDMultihopNodePolicyAttributes(bfdMhNodePol, d)
	if err != nil {
		return nil
	}

	return nil
}
