package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL3outPathAttachment() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL3outPathAttachmentRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"target_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"autostate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encap": &schema.Schema{
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

func dataSourceAciL3outPathAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	tDn := d.Get("target_dn").(string)

	rn := fmt.Sprintf("rspathL3OutAtt-[%s]", tDn)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalInterfaceProfileDn, rn)

	l3extRsPathL3OutAtt, err := getRemoteL3outPathAttachment(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setL3outPathAttachmentAttributes(l3extRsPathL3OutAtt, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
