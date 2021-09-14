package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSpineSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSpineSwitchPolicyGroupRead,
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

func dataSourceAciSpineSwitchPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/funcprof/spaccnodepgrp-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	infraSpineAccNodePGrp, err := getRemoteSpineSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	setSpineSwitchPolicyGroupAttributes(infraSpineAccNodePGrp, d)
	return nil
}
