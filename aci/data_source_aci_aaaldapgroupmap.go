package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciLDAPGroupMap() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciLDAPGroupMapRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ldap",
					"duo",
				}, false),
			},
		})),
	}
}

func dataSourceAciLDAPGroupMapRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	groupType := d.Get("type").(string)

	rn := fmt.Sprintf("userext/%sext/ldapgroupmap-%s", groupType, name)
	dn := fmt.Sprintf("uni/%s", rn)
	aaaLdapGroupMap, err := getRemoteLDAPGroupMap(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setLDAPGroupMapAttributes(aaaLdapGroupMap, d)
	return nil
}
