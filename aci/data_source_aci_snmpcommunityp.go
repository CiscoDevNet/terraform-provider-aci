package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSNMPCommunity() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSNMPCommunityReadContext,
		SchemaVersion: 1,
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

func dataSourceAciSNMPCommunityReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	VRFSNMPCtxDn := d.Get("vrf_snmp_context_dn").(string)
	rn := fmt.Sprintf("community-%s", name)
	dn := fmt.Sprintf("%s/%s", VRFSNMPCtxDn, rn)
	snmpCommunityP, err := getRemoteSNMPCommunity(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setSNMPCommunityAttributes(snmpCommunityP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
