package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTACACSPlusProviderGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciTACACSPlusProviderGroupRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciTACACSPlusProviderGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("userext/tacacsext/tacacsplusprovidergroup-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	aaaTacacsPlusProviderGroup, err := getRemoteTACACSPlusProviderGroup(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setTACACSPlusProviderGroupAttributes(aaaTacacsPlusProviderGroup, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
