package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciDHCPOption() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciDHCPOptionRead,

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
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAciDHCPOptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("opt-%s", name)
	DHCPOptionPolicyDn := d.Get("dhcp_option_policy_dn").(string)

	dn := fmt.Sprintf("%s/%s", DHCPOptionPolicyDn, rn)

	dhcpOption, err := getRemoteDHCPOption(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setDHCPOptionAttributes(dhcpOption, d)
	if err != nil {
		return diag.FromErr(err)
	}
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

func setDHCPOptionAttributes(dhcpOption *models.DHCPOption, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(dhcpOption.DistinguishedName)
	//d.Set("description", dhcpOption.Description)

	if dn != dhcpOption.DistinguishedName {
		d.Set("dhcp_option_policy_dn", "")
	}
	dhcpOptionMap, err := dhcpOption.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", dhcpOptionMap["name"])

	d.Set("annotation", dhcpOptionMap["annotation"])
	d.Set("data", dhcpOptionMap["data"])
	d.Set("dhcp_option_id", dhcpOptionMap["id"])
	d.Set("name_alias", dhcpOptionMap["nameAlias"])
	return d, nil
}
