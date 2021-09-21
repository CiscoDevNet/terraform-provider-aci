package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciManagedNodesZone() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciManagedNodesZoneRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"managed_node_connectivity_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in_band",
					"out_of_band",
				}, false),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciManagedNodesZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("type").(string) == "in_band" {
		return dataSourceAciInBManagedNodesZoneRead(ctx, d, m)
	}
	return dataSourceAciOOBManagedNodesZoneRead(ctx, d, m)
}

func dataSourceAciInBManagedNodesZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	ManagedNodeConnectivityGroupDn := d.Get("managed_node_connectivity_group_dn").(string)
	rn := fmt.Sprintf("inbzone")
	dn := fmt.Sprintf("%s/%s", ManagedNodeConnectivityGroupDn, rn)
	mgmtInBZone, err := getRemoteInBManagedNodesZone(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setInBManagedNodesZoneAttributes(mgmtInBZone, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func dataSourceAciOOBManagedNodesZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	ManagedNodeConnectivityGroupDn := d.Get("managed_node_connectivity_group_dn").(string)
	rn := fmt.Sprintf("oobzone")
	dn := fmt.Sprintf("%s/%s", ManagedNodeConnectivityGroupDn, rn)
	mgmtOOBZone, err := getRemoteOOBManagedNodesZone(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setOOBManagedNodesZoneAttributes(mgmtOOBZone, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
