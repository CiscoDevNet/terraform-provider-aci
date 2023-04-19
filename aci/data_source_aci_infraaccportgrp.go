package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLeafAccessPortPolicyGroup() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciLeafAccessPortPolicyGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciLeafAccessPortPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/funcprof/accportgrp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraAccPortGrp, err := getRemoteLeafAccessPortPolicyGroup(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setLeafAccessPortPolicyGroupAttributes(infraAccPortGrp, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
