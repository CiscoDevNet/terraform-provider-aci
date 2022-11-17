package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSpineInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciSpineInterfaceProfileRead,

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

func dataSourceAciSpineInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/spaccportprof-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraSpAccPortP, err := getRemoteSpineInterfaceProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setSpineInterfaceProfileAttributes(infraSpAccPortP, d)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
