package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciRouteControlProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciRouteControlProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
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

			"route_control_profile_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciRouteControlProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("prof-%s", name)
	ParentDn := d.Get("parent_dn").(string)

	dn := fmt.Sprintf("%s/%s", ParentDn, rn)

	rtctrlProfile, err := getRemoteRouteControlProfile(aciClient, dn)

	if err != nil {
		return err
	}

	d.SetId(dn)
	setRouteControlProfileAttributes(rtctrlProfile, d)
	return nil
}
