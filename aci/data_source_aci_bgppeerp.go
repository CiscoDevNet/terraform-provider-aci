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
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

	rn := fmt.Sprintf("peerP-[%s]", addr)
	ParentDn := d.Get("parent_dn").(string)

	dn := fmt.Sprintf("%s/%s", ParentDn, rn)

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
