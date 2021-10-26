package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPODMaintenanceGroup() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciPODMaintenanceGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"fwtype": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pod_maintenance_group_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciPODMaintenanceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/maintgrp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	maintMaintGrp, err := getRemotePODMaintenanceGroup(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setPODMaintenanceGroupAttributes(maintMaintGrp, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
