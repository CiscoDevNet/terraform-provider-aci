package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSecurityDomain() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciSecurityDomainRead,

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

func dataSourceAciSecurityDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("userext/domain-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	aaaDomain, err := getRemoteSecurityDomain(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setSecurityDomainAttributes(aaaDomain, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
