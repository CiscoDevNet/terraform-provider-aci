package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAAAAuthentication() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciAAAAuthenticationRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"def_role_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ping_check": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"retries": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciAAAAuthenticationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	rnrealm := fmt.Sprintf("userext/authrealm")
	dnrealm := fmt.Sprintf("uni/%s", rnrealm)
	aaaAuthRealm, err := getRemoteAAAAuthentication(aciClient, dnrealm)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dnrealm)
	_, err = setAAAAuthenticationAttributes(aaaAuthRealm, d)
	if err != nil {
		return diag.FromErr(err)
	}
	rnpingep := fmt.Sprintf("userext/pingext")
	dnpingep := fmt.Sprintf("uni/%s", rnpingep)
	aaaPingEp, err := getRemoteDefaultRadiusAuthenticationSettings(aciClient, dnpingep)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setDefaultRadiusAuthenticationSettingsAttributes(aaaPingEp, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
