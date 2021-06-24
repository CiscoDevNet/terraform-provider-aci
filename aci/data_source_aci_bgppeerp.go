package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciBgpPeerConnectivityProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciBgpPeerConnectivityProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_node_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"addr_t_ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"allowed_self_as_cnt": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"peer_ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"private_a_sctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"weight": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"as_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"local_asn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"local_asn_propagate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciBgpPeerConnectivityProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	addr := d.Get("addr").(string)

	rn := fmt.Sprintf("peerP-[%s]", addr)
	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalNodeProfileDn, rn)

	bgpPeerP, err := getRemoteBgpPeerConnectivityProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)
	_, err = setBgpPeerConnectivityProfileAttributes(bgpPeerP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	bgpAsP, err := getRemoteBgpAutonomousSystemProfileFromBgpPeerConnectivityProfile(aciClient, fmt.Sprintf("%s/as", dn))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setBgpAutonomousSystemProfileAttributesFromBgpPeerConnectivityProfile(bgpAsP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	bgpLocalAsnP, err := getRemoteLocalAutonomousSystemProfileFromBgpPeerConnectivityProfile(aciClient, fmt.Sprintf("%s/localasn", dn))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setLocalAutonomousSystemProfileAttributesFromBgpPeerConnectivityProfile(bgpLocalAsnP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
