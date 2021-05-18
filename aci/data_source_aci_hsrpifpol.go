package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciHSRPInterfacePolicy() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciHSRPInterfacePolicyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
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

			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"reload_delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciHSRPInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("hsrpIfPol-%s", name)
	TenantDn := d.Get("tenant_dn").(string)

	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	hsrpIfPol, err := getRemoteHSRPInterfacePolicy(aciClient, dn)

	if err != nil {
		return err
	}

	d.SetId(dn)
	setHSRPInterfacePolicyAttributes(hsrpIfPol, d)
	return nil
}
