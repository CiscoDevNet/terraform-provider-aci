package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLDAPGroupMapruleref() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciLDAPGroupMaprulerefRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"ldap_group_map_dn": &schema.Schema{
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

func dataSourceAciLDAPGroupMaprulerefRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	LDAPGroupMapDn := d.Get("ldap_group_map_dn").(string)
	rn := fmt.Sprintf("ldapgroupmapruleref-%s", name)
	dn := fmt.Sprintf("%s/%s", LDAPGroupMapDn, rn)
	aaaLdapGroupMapRuleRef, err := getRemoteLDAPGroupMapruleref(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setLDAPGroupMaprulerefAttributes(aaaLdapGroupMapRuleRef, d)
	return nil
}
