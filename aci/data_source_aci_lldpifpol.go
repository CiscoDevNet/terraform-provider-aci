package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLLDPInterfacePolicy() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciLLDPInterfacePolicyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"admin_rx_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"admin_tx_st": &schema.Schema{
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

func dataSourceAciLLDPInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/lldpIfP-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	lldpIfPol, err := getRemoteLLDPInterfacePolicy(aciClient, dn)

	if err != nil {
		return err
	}
	setLLDPInterfacePolicyAttributes(lldpIfPol, d)
	return nil
}
