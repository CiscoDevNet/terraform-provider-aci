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

func dataSourceAciIGMPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciIGMPInterfacePolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"last_member_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_member_response_time": {
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
			"query_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"robustness_variable": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"response_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"startup_query_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"startup_query_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maximum_mulitcast_entries": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reserved_mulitcast_entries": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_limit_route_map": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"report_policy_route_map": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"static_report_route_map": {
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

	_, err = getandSetIGMPIfPolRelationshipAttributes(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] IGMP Interface Policy Relationship Attributes - Read finished successfully")
	}

	return nil
}
