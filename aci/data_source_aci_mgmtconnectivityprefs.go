package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciMgmtconnectivitypreference() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciMgmtconnectivitypreferenceRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"interface_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciMgmtconnectivitypreferenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("fabric/connectivityPrefs")
	dn := fmt.Sprintf("uni/%s", rn)
	mgmtConnectivityPrefs, err := getRemoteMgmtconnectivitypreference(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setMgmtconnectivitypreferenceAttributes(mgmtConnectivityPrefs, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
