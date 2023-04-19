package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciUserRole() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciUserRoleReadContext,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"user_domain_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"priv_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciUserRoleReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	UserDomainDn := d.Get("user_domain_dn").(string)
	rn := fmt.Sprintf("role-%s", name)
	dn := fmt.Sprintf("%s/%s", UserDomainDn, rn)
	aaaUserRole, err := getRemoteUserRole(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setUserRoleAttributes(aaaUserRole, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
