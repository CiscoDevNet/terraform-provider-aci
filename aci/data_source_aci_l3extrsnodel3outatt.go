package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciFabricNode() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciFabricNodeRead,

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

func dataSourceAciFabricNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	tdn := d.Get("tdn").(string)

	rn := fmt.Sprintf("rsnodeL3OutAtt-[%s]", tdn)
	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalNodeProfileDn, rn)

	l3extRsNodeL3OutAtt, err := getRemoteFabricNode(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setFabricNodeAttributes(l3extRsNodeL3OutAtt, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
