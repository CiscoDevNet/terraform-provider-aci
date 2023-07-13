package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciPowerSupplyRedundancyPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciPowerSupplyRedundancyPolicyRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"admin_rdn_m": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciPowerSupplyRedundancyPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/psuInstP-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)

	psuInstPol, err := getRemotePowerSupplyRedundancyPolicy(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setPowerSupplyRedundancyPolicyAttributes(psuInstPol, d)
	if err != nil {
		return nil
	}

	return nil
}
