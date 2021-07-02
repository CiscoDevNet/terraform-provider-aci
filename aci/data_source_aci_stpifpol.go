package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSpanningTreeInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSpanningTreeInterfacePolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ctrl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciSpanningTreeInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/ifPol-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	stpIfPol, err := getRemoteSpanningTreeInterfacePolicy(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setSpanningTreeInterfacePolicyAttributes(stpIfPol, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
