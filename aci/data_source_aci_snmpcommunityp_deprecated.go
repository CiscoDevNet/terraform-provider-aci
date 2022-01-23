package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSNMPCommunityDeprecated() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use aci_snmp_community data source instead",
		ReadContext:        dataSourceAciSNMPCommunityReadDeprecated,
		SchemaVersion:      1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_snmp_context_dn": &schema.Schema{
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

func dataSourceAciSNMPCommunityReadDeprecated(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	parentDn := d.Get("vrf_snmp_context_dn").(string)
	dn := fmt.Sprintf(models.DnsnmpCommunityP, parentDn, name)

	snmpCommunityP, err := getRemoteSNMPCommunity(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setSNMPCommunityAttributesDeprecated(snmpCommunityP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
