package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciL3outStaticRouteNextHop() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciL3outStaticRouteNextHopRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"static_route_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"nh_addr": &schema.Schema{
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

			"pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"nexthop_profile_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciL3outStaticRouteNextHopRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	nhAddr := d.Get("nh_addr").(string)

	rn := fmt.Sprintf("nh-[%s]", nhAddr)
	StaticRouteDn := d.Get("static_route_dn").(string)

	dn := fmt.Sprintf("%s/%s", StaticRouteDn, rn)

	ipNexthopP, err := getRemoteL3outStaticRouteNextHop(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setL3outStaticRouteNextHopAttributes(ipNexthopP, d)
	return nil
}
