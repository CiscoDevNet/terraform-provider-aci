package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTACACSMonitoringDestinationGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciTACACSMonitoringDestinationGroupReadContext,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciTACACSMonitoringDestinationGroupReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/tacacsgroup-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	tacacsGroup, err := getRemoteTACACSMonitoringDestinationGroup(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setTACACSMonitoringDestinationGroupAttributes(tacacsGroup, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
