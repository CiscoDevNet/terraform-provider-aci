package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciVSwitchPolicyGroupRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciVSwitchPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)
	rn := fmt.Sprintf("vswitchpolcont")
	dn := fmt.Sprintf("%s/%s", VMMDomainDn, rn)
	vmmVSwitchPolicyCont, err := getRemoteVSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setVSwitchPolicyGroupAttributes(vmmVSwitchPolicyCont, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
