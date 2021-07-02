package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciNodeBlockMG() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciNodeBlockReadMG,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"pod_maintenance_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"from_": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"to_": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciNodeBlockReadMG(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("nodeblk-%s", name)
	PODMaintenanceGroupDn := d.Get("pod_maintenance_group_dn").(string)

	dn := fmt.Sprintf("%s/%s", PODMaintenanceGroupDn, rn)

	fabricNodeBlk, err := getRemoteNodeBlockMG(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setNodeBlockAttributesMG(fabricNodeBlk, d)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
