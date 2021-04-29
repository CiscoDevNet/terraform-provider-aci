package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciDHCPOption() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciDHCPOptionRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
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
				Computed: true,
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
		}),
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

func getRemoteDHCPOption(client *client.Client, dn string) (*models.DHCPOption, error) {
	dhcpOptionCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpOption := models.DHCPOptionFromContainer(dhcpOptionCont)

	if dhcpOption.DistinguishedName == "" {
		return nil, fmt.Errorf("DHCPOption %s not found", dhcpOption.DistinguishedName)
	}

	return dhcpOption, nil
}

func setDHCPOptionAttributes(dhcpOption *models.DHCPOption, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(dhcpOption.DistinguishedName)
	//d.Set("description", dhcpOption.Description)

	if dn != dhcpOption.DistinguishedName {
		d.Set("dhcp_option_policy_dn", "")
	}
	dhcpOptionMap, _ := dhcpOption.ToMap()

	d.Set("name", dhcpOptionMap["name"])

	d.Set("annotation", dhcpOptionMap["annotation"])
	d.Set("data", dhcpOptionMap["data"])
	d.Set("dhcp_option_id", dhcpOptionMap["id"])
	d.Set("name_alias", dhcpOptionMap["nameAlias"])
	return d
}
