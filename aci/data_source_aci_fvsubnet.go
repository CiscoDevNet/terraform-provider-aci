package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSubnet() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciSubnetRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"preferred": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"virtual": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciSubnetRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	ip := d.Get("ip").(string)

	rn := fmt.Sprintf("subnet-[%s]", ip)
	BridgeDomainDn := d.Get("parent_dn").(string)

	dn := fmt.Sprintf("%s/%s", BridgeDomainDn, rn)

	fvSubnet, err := getRemoteSubnet(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setSubnetAttributes(fvSubnet, d)
	return nil
}
