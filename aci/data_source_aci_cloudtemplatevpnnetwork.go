package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudTemplateforVPNNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciCloudTemplateforVPNNetworkRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"aci_cloud_external_network_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_site_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"remote_site_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipsec_tunnel": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike_version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"public_ip_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet_pool_dn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"pre_shared_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"bgp_peer_asn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		})),
	}
}

func dataSourceAciCloudTemplateforVPNNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TemplateforExternalNetworkDn := d.Get("aci_cloud_external_network_dn").(string)
	rn := fmt.Sprintf("vpnnetwork-%s", name)
	dn := fmt.Sprintf("%s/%s", TemplateforExternalNetworkDn, rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)

	cloudtemplateVpnNetwork, err := getRemoteTemplateforVPNNetwork(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setTemplateforVPNNetworkAttributes(cloudtemplateVpnNetwork, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Data Source - Begining Read of cloud IPSec Tunnel attributes.")
	cloudtemplateIpSecTunnelData, err := aciClient.ListCloudTemplateforIpSectunnel(dn)
	if err != nil {
		log.Printf("[DEBUG] Data Source - Error while reading cloud IPSec Tunnel attributes %v", err)
	}

	cloudtemplateIpSecTunnelSet := make([]map[string]string, 0, 1)
	for _, cloudtemplateIpSecTunnel := range cloudtemplateIpSecTunnelData {

		cloudIpSecTunnelAttMap, cloudtemplateIpSecTunnelDn, err := setCloudTemplateforIpSecTunnelAttributes(cloudtemplateIpSecTunnel, make(map[string]string))
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Data Source - Begining Read of cloud BGP IPV4 Peer attributes.")
		bgpIPv4PeerData, err := aciClient.ListBGPIPv4Peer(cloudtemplateIpSecTunnelDn)
		if err != nil {
			log.Printf("[DEBUG] Data Source - Error while reading cloud BGP IPV4 Peer attributes %v", err)
		}
		for _, bgpIPv4Peer := range bgpIPv4PeerData {
			bgpPeerAsnAtt, err := getASNfromBGPTPV4Peer(bgpIPv4Peer, make(map[string]string))
			if err != nil {
				d.SetId("")
				return diag.FromErr(err)
			}
			cloudIpSecTunnelAttMap["bgp_peer_asn"] = bgpPeerAsnAtt["bgp_peer_asn_att"]
		}
		log.Printf("[DEBUG] Data Source - Read cloud BGP IPV4 Peer finished successfully.")

		cloudtemplateIpSecTunnelSet = append(cloudtemplateIpSecTunnelSet, cloudIpSecTunnelAttMap)
		log.Printf("LOGs Data Source - cloudtemplateIpSecTunnelSet. : %v ", cloudtemplateIpSecTunnelSet)
	}
	d.Set("ipsec_tunnel", cloudtemplateIpSecTunnelSet)
	log.Printf("LOGs Data Source - d. : %v ", d)
	log.Printf("[DEBUG] Data Source - Read cloud IPSec Tunnel finished successfully.")

	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}
