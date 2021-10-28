package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSAMLProviderGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSAMLProviderGroupReadContext,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciSAMLProviderGroupReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("userext/samlext/samlprovidergroup-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	aaaSamlProviderGroup, err := getRemoteSAMLProviderGroup(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setSAMLProviderGroupAttributes(aaaSamlProviderGroup, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
