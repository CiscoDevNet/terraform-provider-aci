package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciBfdMultihopInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciBfdMultihopInterfacePolicyRead,
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
			"min_receive_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_transmit_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"descr": {
				Type:     schema.TypeString,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciBfdMultihopInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnbfdMhIfPol, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	bfdMhIfPol, err := getAciBfdMultihopInterfacePolicy(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setAciBfdMultihopInterfacePolicyAttributes(bfdMhIfPol, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
