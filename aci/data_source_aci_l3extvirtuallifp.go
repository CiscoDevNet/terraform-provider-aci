package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
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
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"autostate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encap_scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"if_inst_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv6_dad": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ll_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mtu": &schema.Schema{
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
	return nil
}
