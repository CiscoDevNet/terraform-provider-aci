package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciMgmtStaticNode() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciAciMgmtStaticNodeRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"management_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"t_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

			"gw": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"v6_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"v6_gw": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciAciMgmtStaticNodeRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	tDn := d.Get("t_dn").(string)

	rn := fmt.Sprintf("rsooBStNode-[%s]", tDn)
	managementEPgDn := d.Get("management_epg_dn").(string)

	dn := fmt.Sprintf("%s/%s", managementEPgDn, rn)

	bandType := d.Get("type").(string)

	if bandType == "in-band" {
		mgmtRsInBStNode, err := getRemoteInbandStaticNode(aciClient, dn)
		if err != nil {
			d.SetId("")
			return nil
		}

		d.SetId(dn)
		setMgmtStaticNodeAttributes(nil, mgmtRsInBStNode, "in-band", d)

	} else {
		mgmtRsOoBStNode, err := getRemoteOutofbandStaticNode(aciClient, dn)
		if err != nil {
			d.SetId("")
			return nil
		}

		d.SetId(dn)
		setMgmtStaticNodeAttributes(mgmtRsOoBStNode, nil, "out-of-band", d)
	}
	return nil
}
