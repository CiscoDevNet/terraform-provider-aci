package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLeafBreakoutPortGroup() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciLeafBreakoutPortGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"brkout_map": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciLeafBreakoutPortGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/funcprof/brkoutportgrp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraBrkoutPortGrp, err := getRemoteLeafBreakoutPortGroup(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setLeafBreakoutPortGroupAttributes(infraBrkoutPortGrp, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
