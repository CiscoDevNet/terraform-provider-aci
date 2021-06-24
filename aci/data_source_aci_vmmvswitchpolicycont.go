package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciVSwitchPolicyGroupRead,
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

func dataSourceAciVSwitchPolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)
	rn := fmt.Sprintf("vswitchpolcont")
	dn := fmt.Sprintf("%s/%s", VMMDomainDn, rn)
	vmmVSwitchPolicyCont, err := getRemoteVSwitchPolicyGroup(aciClient, dn)
	if err != nil {
		return err
	}
	setVSwitchPolicyGroupAttributes(vmmVSwitchPolicyCont, d)
	return nil
}
