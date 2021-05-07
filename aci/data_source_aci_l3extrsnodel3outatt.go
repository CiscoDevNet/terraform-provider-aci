package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciFabricNode() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciFabricNodeRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_node_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"config_issues": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rtr_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rtr_id_loop_back": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciFabricNodeRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	rn := fmt.Sprintf("rsnodeL3OutAtt-[%s]", tDn)
	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalNodeProfileDn, rn)

	l3extRsNodeL3OutAtt, err := getRemoteFabricNode(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setFabricNodeAttributes(l3extRsNodeL3OutAtt, d)
	return nil
}
