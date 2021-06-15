package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVXLANPool() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciVXLANPoolRead,

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

func dataSourceAciVXLANPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/vxlanns-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	fvnsVxlanInstP, err := getRemoteVXLANPool(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	setVXLANPoolAttributes(fvnsVxlanInstP, d)
	return nil
}
