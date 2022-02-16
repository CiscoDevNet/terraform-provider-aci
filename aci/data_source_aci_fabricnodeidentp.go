package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciFabricNodeMember() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciFabricNodeMemberRead,

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

			"name": &schema.Schema{
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

func dataSourceAciFabricNodeMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	serial := d.Get("serial").(string)

	rn := fmt.Sprintf("controller/nodeidentpol/nodep-%s", serial)

	dn := fmt.Sprintf("uni/%s", rn)

	fabricNodeIdentP, err := getRemoteFabricNodeMember(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setFabricNodeMemberAttributes(fabricNodeIdentP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
