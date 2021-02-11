package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciDHCPOption() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciDHCPOptionRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"dhcp_option_policy_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "orchestrator:terraform",
			},

			"data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dhcp_option_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAciDHCPOptionRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("opt-%s", name)
	DHCPOptionPolicyDn := d.Get("dhcp_option_policy_dn").(string)

	dn := fmt.Sprintf("%s/%s", DHCPOptionPolicyDn, rn)

	dhcpOption, err := getRemoteDHCPOption(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setDHCPOptionAttributes(dhcpOption, d)
	return nil
}
