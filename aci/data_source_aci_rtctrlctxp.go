package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciRouteControlContext() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciRouteControlContextRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"route_control_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"set_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to rtctrl:AttrP",
			},
		})),
	}
}

func dataSourceAciRouteControlContextRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	RouteControlProfileDn := d.Get("route_control_profile_dn").(string)
	rn := fmt.Sprintf("ctx-%s", name)
	dn := fmt.Sprintf("%s/%s", RouteControlProfileDn, rn)
	rtctrlCtxP, err := getRemoteRouteControlContext(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setRouteControlContextAttributes(rtctrlCtxP, d)
	return nil
}
