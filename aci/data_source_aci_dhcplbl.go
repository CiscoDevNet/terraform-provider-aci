package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciDHCPRelayLabel() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciDHCPRelayLabelRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"bridge_domain_dn": &schema.Schema{
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

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciDHCPRelayLabelRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("dhcplbl-%s", name)
	BridgeDomainDn := d.Get("bridge_domain_dn").(string)

	dn := fmt.Sprintf("%s/%s", BridgeDomainDn, rn)

	dhcpLbl, err := getRemoteDHCPRelayLabel(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setDHCPRelayLabelAttributes(dhcpLbl, d)
	return nil
}
