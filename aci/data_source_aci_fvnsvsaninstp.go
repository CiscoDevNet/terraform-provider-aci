package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVSANPool() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciVSANPoolRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"alloc_mode": &schema.Schema{
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

func dataSourceAciVSANPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	allocMode := d.Get("alloc_mode").(string)

	rn := fmt.Sprintf("infra/vsanns-[%s]-%s", name, allocMode)

	dn := fmt.Sprintf("uni/%s", rn)

	fvnsVsanInstP, err := getRemoteVSANPool(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	setVSANPoolAttributes(fvnsVsanInstP, d)
	return nil
}
