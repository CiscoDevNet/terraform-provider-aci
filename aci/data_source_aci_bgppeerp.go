package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciBgpPeerConnectivityProfile() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciBgpPeerConnectivityProfileRead,

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

func dataSourceAciBgpPeerConnectivityProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	addr := d.Get("addr").(string)

	rn := fmt.Sprintf("peerP-[%s]", addr)
	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalNodeProfileDn, rn)

	bgpPeerP, err := getRemoteBgpPeerConnectivityProfile(aciClient, dn)

	if err != nil {
		return err
	}

	d.SetId(dn)
	setBgpPeerConnectivityProfileAttributes(bgpPeerP, d)

	bgpAsP, err := getRemoteBgpAutonomousSystemProfileFromBgpPeerConnectivityProfile(aciClient, fmt.Sprintf("%s/as", dn))
	if err != nil {
		return err
	}
	setBgpAutonomousSystemProfileAttributesFromBgpPeerConnectivityProfile(bgpAsP, d)

	bgpLocalAsnP, err := getRemoteLocalAutonomousSystemProfileFromBgpPeerConnectivityProfile(aciClient, fmt.Sprintf("%s/localasn", dn))
	if err != nil {
		return err
	}
	setLocalAutonomousSystemProfileAttributesFromBgpPeerConnectivityProfile(bgpLocalAsnP, d)

	return nil
}
