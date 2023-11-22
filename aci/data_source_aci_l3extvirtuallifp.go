package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVirtualLogicalInterfaceProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciVirtualLogicalInterfaceProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"node_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"autostate": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"encap_scope": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"if_inst_t": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv6_dad": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ll_addr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"relation_l3ext_rs_dyn_path_att": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"floating_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forged_transmit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_change": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"promiscuous_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enhanced_lag_policy_dn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encap": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func dataSourceAciVirtualLogicalInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	nodeDn := d.Get("node_dn").(string)

	encap := d.Get("encap").(string)

	rn := fmt.Sprintf("vlifp-[%s]-[%s]", nodeDn, encap)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalInterfaceProfileDn, rn)

	l3extVirtualLIfP, err := getRemoteVirtualLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)
	_, err = setVirtualLogicalInterfaceProfileAttributes(l3extVirtualLIfP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getAndSetL3extRsDynPathAttFromLogicalInterfaceProfile(aciClient, dn, d)

	return nil
}
