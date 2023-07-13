package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPIMInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciPIMInterfacePolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"control_state": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"designated_router_delay": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"designated_router_priority": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hello_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"join_prune_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inbound_join_prune_filter_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"outbound_join_prune_filter_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"neighbor_filter_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciPIMInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnPimIfPol, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	pimIfPol, err := getRemotePIMInterfacePolicy(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setPIMInterfacePolicyAttributes(pimIfPol, d)
	if err != nil {
		return nil
	}

	_, err = getandSetPIMIfPolRelationshipAttributes(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] PimIfPol Relationship Attributes - Read finished successfully")
	}

	return nil
}
