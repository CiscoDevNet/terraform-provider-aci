package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciIPSLAMonitoringPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciIPSLAMonitoringPolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type_of_service": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"traffic_class_value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_data_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sla_detect_multiplier": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sla_frequency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sla_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sla_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciIPSLAMonitoringPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnfvIPSLAMonitoringPol, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	fvIPSLAMonitoringPol, err := getRemoteIPSLAMonitoringPolicy(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setIPSLAMonitoringPolicyAttributes(fvIPSLAMonitoringPol, d)
	if err != nil {
		return nil
	}

	return nil
}
