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
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use parent_dn instead",
			},

			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"addr_t_ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"private_a_sctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	var parentDn string

	rn := fmt.Sprintf("peerP-[%s]", addr)
	if d.Get("logical_node_profile_dn").(string) != "" && d.Get("parent_dn").(string) != "" {
		return diag.FromErr(fmt.Errorf("Usage of both parent_dn and logical_node_profile_dn parameters is not supported. logical_node_profile_dn parameter will be deprecated use parent_dn instead."))
	} else if d.Get("parent_dn").(string) != "" {
		parentDn = d.Get("parent_dn").(string)
	} else if d.Get("logical_node_profile_dn").(string) != "" {
		parentDn = d.Get("logical_node_profile_dn").(string)
	} else {
		return diag.FromErr(fmt.Errorf("parent_dn is required to query a BGP Peer Connectivity Profile"))
	}

	dn := fmt.Sprintf("%s/%s", parentDn, rn)

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
