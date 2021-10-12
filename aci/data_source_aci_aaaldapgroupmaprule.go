package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciLDAPGroupMapRule() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciLDAPGroupMapRuleReadContext,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"groupdn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"duo", "ldap"}, false),
			},
		})),
	}
}

func dataSourceAciLDAPGroupMapRuleReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	ldap_group_map_rule_type := d.Get("type").(string)
	var rn string
	if ldap_group_map_rule_type == "duo" {
		rn = fmt.Sprintf("userext/duoext/ldapgroupmaprule-%s", name)
	} else {
		rn = fmt.Sprintf("userext/ldapext/ldapgroupmaprule-%s", name)
	}
	dn := fmt.Sprintf("uni/%s", rn)
	aaaLdapGroupMapRule, err := getRemoteLDAPGroupMapRule(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setLDAPGroupMapRuleAttributes(ldap_group_map_rule_type, aaaLdapGroupMapRule, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
