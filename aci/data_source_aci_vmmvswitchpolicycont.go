package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciVSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciVSwitchPolicyGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciVSwitchPolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("vswitchpolcont")
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	dn := fmt.Sprintf("%s/%s", VMMDomainDn, rn)

	vmmVSwitchPolicyCont, err := getRemoteVSwitchPolicyGroup(aciClient, dn)

	if err != nil {
		return err
	}
	setVSwitchPolicyGroupAttributes(vmmVSwitchPolicyCont, d)
	return nil
}
