package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTenant() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciTenantRead,

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

func dataSourceAciTenantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("tn-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	fvTenant, err := getRemoteTenant(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setTenantAttributes(fvTenant, d)
	
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
