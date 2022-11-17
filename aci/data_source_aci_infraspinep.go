package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSpineProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciSpineProfileRead,

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

func dataSourceAciSpineProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/spprof-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	infraSpineP, err := getRemoteSpineProfile(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = setSpineProfileAttributes(infraSpineP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
