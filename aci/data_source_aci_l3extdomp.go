package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3DomainProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL3DomainProfileRead,

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

func dataSourceAciL3DomainProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("l3dom-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	l3extDomP, err := getRemoteL3DomainProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setL3DomainProfileAttributes(l3extDomP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
