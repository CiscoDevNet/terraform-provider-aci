package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciL2outExternalEpg() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciL2outExternalEpgRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l2_outside_dn": &schema.Schema{
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

			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flood_on_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciL2outExternalEpgRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("instP-%s", name)
	L2OutsideDn := d.Get("l2_outside_dn").(string)

	dn := fmt.Sprintf("%s/%s", L2OutsideDn, rn)

	l2extInstP, err := getRemoteL2outExternalEpg(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setL2outExternalEpgAttributes(l2extInstP, d)
	return nil
}
