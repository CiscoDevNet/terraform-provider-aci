package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3outStaticRouteNextHop() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL3outStaticRouteNextHopRead,

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

func dataSourceAciL3outStaticRouteNextHopRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	nhAddr := d.Get("nh_addr").(string)

	rn := fmt.Sprintf("nh-[%s]", nhAddr)
	StaticRouteDn := d.Get("static_route_dn").(string)

	dn := fmt.Sprintf("%s/%s", StaticRouteDn, rn)

	ipNexthopP, err := getRemoteL3outStaticRouteNextHop(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setL3outStaticRouteNextHopAttributes(ipNexthopP, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
