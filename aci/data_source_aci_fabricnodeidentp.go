package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciFabricNodeMember() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciFabricNodeMemberRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"serial": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ext_pool_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fabric_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"node_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"node_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pod_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciFabricNodeMemberRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	serial := d.Get("serial").(string)

	rn := fmt.Sprintf("controller/nodeidentpol/nodep-%s", serial)

	dn := fmt.Sprintf("uni/%s", rn)

	fabricNodeIdentP, err := getRemoteFabricNodeMember(aciClient, dn)

	if err != nil {
		return err
	}
	setFabricNodeMemberAttributes(fabricNodeIdentP, d)
	return nil
}
