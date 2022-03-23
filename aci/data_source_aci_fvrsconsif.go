package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciContractInterfaceRelationship() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciContractInterfaceRelationshipRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"application_epg_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"prio": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"contract_interface_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciContractInterfaceRelationshipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	tnVzCPIfName := GetMOName(d.Get("contract_interface_dn").(string))
	ApplicationEPGDn := d.Get("application_epg_dn").(string)
	rn := fmt.Sprintf("rsconsIf-%s", tnVzCPIfName)
	dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

	fvRsConsIf, err := getRemoteContractInterfaceRelationship(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setContractInterfaceRelationshipAttributes(fvRsConsIf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
