package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPhysicalDomain() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciPhysicalDomainRead,

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

func dataSourceAciPhysicalDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("phys-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	physDomP, err := getRemotePhysicalDomain(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	setPhysicalDomainAttributes(physDomP, d)
	return nil
}
