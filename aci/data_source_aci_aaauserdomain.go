package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciUserDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciUserDomainRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"local_user_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciUserDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	LocalUserDn := d.Get("local_user_dn").(string)
	rn := fmt.Sprintf("userdomain-%s", name)
	dn := fmt.Sprintf("%s/%s", LocalUserDn, rn)
	aaaUserDomain, err := getRemoteUserDomain(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setUserDomainAttributes(aaaUserDomain, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
