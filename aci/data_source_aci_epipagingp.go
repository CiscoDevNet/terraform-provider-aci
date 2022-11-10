package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciIPAgingPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciIPAgingPolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciIPAgingPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := "default"

	rn := fmt.Sprintf("infra/ipAgingP-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	epIpAgingP, err := getRemoteIPAgingPolicy(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setIPAgingPolicyAttributes(epIpAgingP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
