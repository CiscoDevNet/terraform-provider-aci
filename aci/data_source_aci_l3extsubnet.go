package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciL3ExtSubnet() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciL3ExtSubnetRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"external_network_instance_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"aggregate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciL3ExtSubnetRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	ip := d.Get("ip").(string)

	rn := fmt.Sprintf("extsubnet-[%s]", ip)
	ExternalNetworkInstanceProfileDn := d.Get("external_network_instance_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", ExternalNetworkInstanceProfileDn, rn)

	l3extSubnet, err := getRemoteSubnet(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setSubnetAttributes(l3extSubnet, d)
	return nil
}
