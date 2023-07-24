package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSnmpUserProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSnmpUserProfileRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"snmp_policy_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authorization_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"privacy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciSnmpUserProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	snmpPolicyDn := d.Get("snmp_policy_dn").(string)
	rn := fmt.Sprintf(models.RnSnmpUserP, name)
	dn := fmt.Sprintf("%s/%s", snmpPolicyDn, rn)

	snmpUserP, err := getRemoteSnmpUserProfile(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setSnmpUserProfileAttributes(snmpUserP, d)
	if err != nil {
		return nil
	}

	return nil
}
