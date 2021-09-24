package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSNMPContextProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciSNMPContextProfileRead,
		SchemaVersion: 1,
		Schema: AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciSNMPContextProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	VRFDn := d.Get("vrf_dn").(string)
	rn := fmt.Sprintf("snmpctx")
	dn := fmt.Sprintf("%s/%s", VRFDn, rn)
	snmpCtxP, err := getRemoteSNMPContextProfile(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setSNMPContextProfileAttributes(snmpCtxP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
