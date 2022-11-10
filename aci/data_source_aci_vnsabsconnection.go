package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciConnection() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciConnectionRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l4_l7_service_graph_template_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"adj_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"conn_dir": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"conn_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"direct_connect": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"unicast_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciConnectionRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("AbsConnection-%s", name)
	L4L7ServiceGraphTemplateDn := d.Get("l4_l7_service_graph_template_dn").(string)

	dn := fmt.Sprintf("%s/%s", L4L7ServiceGraphTemplateDn, rn)

	vnsAbsConnection, err := getRemoteConnection(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setConnectionAttributes(vnsAbsConnection, d)
	return nil
}
