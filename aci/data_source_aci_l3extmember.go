package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAciL3outVPCMember() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciL3outVPCMemberRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"leaf_port_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"side": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A",
					"B",
				}, false),
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv6_dad": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ll_addr": &schema.Schema{
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

func dataSourceAciL3outVPCMemberRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	side := d.Get("side").(string)

	rn := fmt.Sprintf("mem-%s", side)
	LeafPortDn := d.Get("leaf_port_dn").(string)

	dn := fmt.Sprintf("%s/%s", LeafPortDn, rn)

	l3extMember, err := getRemoteL3outVPCMember(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setL3outVPCMemberAttributes(l3extMember, d)
	return nil
}
