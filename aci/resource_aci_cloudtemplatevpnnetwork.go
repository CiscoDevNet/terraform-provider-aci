package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciCloudTemplateforVPNNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudTemplateforVPNNetworkCreate,
		UpdateContext: resourceAciCloudTemplateforVPNNetworkUpdate,
		ReadContext:   resourceAciCloudTemplateforVPNNetworkRead,
		DeleteContext: resourceAciCloudTemplateforVPNNetworkDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudTemplateforVPNNetworkImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"aci_cloud_external_network_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
							ValidateFunc: validation.StringInSlice([]string{
								"ikev1",
								"ikev2",
							}, false),
						},
						"public_ip_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet_pool_name": {
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

func getRemoteTemplateforVPNNetwork(client *client.Client, dn string) (*models.TemplateforVPNNetwork, error) {
	cloudtemplateVpnNetworkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudtemplateVpnNetwork := models.TemplateforVPNNetworkFromContainer(cloudtemplateVpnNetworkCont)
	if cloudtemplateVpnNetwork.DistinguishedName == "" {
		return nil, fmt.Errorf("TemplateforVPNNetwork %s not found", cloudtemplateVpnNetwork.DistinguishedName)
	}
	return cloudtemplateVpnNetwork, nil
}

func setTemplateforVPNNetworkAttributes(cloudtemplateVpnNetwork *models.TemplateforVPNNetwork, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudtemplateVpnNetwork.DistinguishedName)

	cloudtemplateVpnNetworkMap, err := cloudtemplateVpnNetwork.ToMap()
	if err != nil {
		return d, err
	}

	if dn != cloudtemplateVpnNetwork.DistinguishedName {
		d.Set("aci_cloud_external_network_dn", "")
	} else {
		d.Set("aci_cloud_external_network_dn", GetParentDn(dn, "/"+fmt.Sprintf(models.RncloudtemplateVpnNetwork, cloudtemplateVpnNetworkMap["name"])))
	}
	d.Set("annotation", cloudtemplateVpnNetworkMap["annotation"])
	d.Set("name", cloudtemplateVpnNetworkMap["name"])
	d.Set("remote_site_id", cloudtemplateVpnNetworkMap["remoteSiteId"])
	d.Set("remote_site_name", cloudtemplateVpnNetworkMap["remoteSiteName"])
	d.Set("name_alias", cloudtemplateVpnNetworkMap["nameAlias"])

	return d, nil
}

func setCloudTemplateforIpSecTunnelAttributes(cloudtemplateIpSecTunnel *models.CloudTemplateforIpSectunnel, d map[string]string) (map[string]string, string, error) {
	cloudtemplateIpSecTunnelMap, err := cloudtemplateIpSecTunnel.ToMap()
	if err != nil {
		return nil, "", err
	}

	d = map[string]string{
		"ike_version":       cloudtemplateIpSecTunnelMap["ikeVersion"],
		"public_ip_address": cloudtemplateIpSecTunnelMap["peeraddr"],
		"subnet_pool_name":  cloudtemplateIpSecTunnelMap["poolname"],
		"pre_shared_key":    cloudtemplateIpSecTunnelMap["preSharedKey"],
	}

	return d, cloudtemplateIpSecTunnel.DistinguishedName, nil
}

func getASNfromBGPTPV4Peer(cloudtemplateBgpIpv4 *models.BGPIPv4Peer, d map[string]string) (map[string]string, error) {

	cloudtemplateBgpIpv4Map, err := cloudtemplateBgpIpv4.ToMap()
	if err != nil {
		return d, err
	}

	log.Printf("LOGs getASNfromBGPTPV4Peer ASN : %v ", cloudtemplateBgpIpv4Map["peerasn"])
	d = map[string]string{
		"bgp_peer_asn_att": cloudtemplateBgpIpv4Map["peerasn"],
	}
	return d, nil
}

func resourceAciCloudTemplateforVPNNetworkImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudtemplateVpnNetwork, err := getRemoteTemplateforVPNNetwork(aciClient, dn)
	if err != nil {
		return nil, err
	}

	schemaFilled, err := setTemplateforVPNNetworkAttributes(cloudtemplateVpnNetwork, d)
	if err != nil {
		return nil, err
	}

	cloudtemplateIpSecTunnelData, err := aciClient.ListCloudTemplateforIpSectunnel(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while importing cloud IPSec Tunnel attributes %v", err)
	}
	cloudtemplateIpSecTunnelSet := make([]map[string]string, 0, 1)
	for _, cloudtemplateIpSecTunnel := range cloudtemplateIpSecTunnelData {

		cloudIpSecTunnelAttMap, cloudtemplateIpSecTunnelDn, err := setCloudTemplateforIpSecTunnelAttributes(cloudtemplateIpSecTunnel, make(map[string]string))
		if err != nil {
			d.SetId("")
			return nil, err
		}
		log.Printf("LOGs IMPORT IPSEC cloudIpSecTunnelAttMap : %v    nand TYPE : %v", cloudIpSecTunnelAttMap, reflect.TypeOf(cloudIpSecTunnelAttMap))

		bgpIPv4PeerData, err := aciClient.ListBGPIPv4Peer(cloudtemplateIpSecTunnelDn)
		if err != nil {
			log.Printf("[DEBUG] Error while reading cloud IPSec Tunnel attributes %v", err)
		}
		for _, bgpIPv4Peer := range bgpIPv4PeerData {
			log.Printf("LOGs IMPORT FOR bgpIPv4PeerData : %v ", bgpIPv4Peer)
			bgpPeerAsnAtt, err := getASNfromBGPTPV4Peer(bgpIPv4Peer, make(map[string]string))
			if err != nil {
				d.SetId("")
				return nil, err
			}
			log.Printf("LOGs IMPORT bgpPeerAsnAtt ASN : %v ", bgpPeerAsnAtt["bgp_peer_asn_att"])
			cloudIpSecTunnelAttMap["bgp_peer_asn"] = bgpPeerAsnAtt["bgp_peer_asn_att"]
		}

		cloudtemplateIpSecTunnelSet = append(cloudtemplateIpSecTunnelSet, cloudIpSecTunnelAttMap)
		log.Printf("LOGs IMPORT IP Tunnel Set : %v ", cloudtemplateIpSecTunnelSet)
	}
	d.Set("ipsec_tunnel", cloudtemplateIpSecTunnelSet)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func TypeOf(cloudIpSecTunnelAttMap map[string]string) {
	panic("unimplemented")
}

func resourceAciCloudTemplateforVPNNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TemplateforVPNNetwork: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TemplateforExternalNetworkDn := d.Get("aci_cloud_external_network_dn").(string)

	cloudtemplateVpnNetworkAttr := models.TemplateforVPNNetworkAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateVpnNetworkAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateVpnNetworkAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudtemplateVpnNetworkAttr.Name = Name.(string)
	}

	if RemoteSiteId, ok := d.GetOk("remote_site_id"); ok {
		cloudtemplateVpnNetworkAttr.RemoteSiteId = RemoteSiteId.(string)
	}

	if RemoteSiteName, ok := d.GetOk("remote_site_name"); ok {
		cloudtemplateVpnNetworkAttr.RemoteSiteName = RemoteSiteName.(string)
	}
	cloudtemplateVpnNetwork := models.NewTemplateforVPNNetwork(fmt.Sprintf(models.RncloudtemplateVpnNetwork, name), TemplateforExternalNetworkDn, nameAlias, cloudtemplateVpnNetworkAttr)

	err := aciClient.Save(cloudtemplateVpnNetwork)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("LOGs CREATE cloudtemplateVpnNetwork : %v ", cloudtemplateVpnNetwork)
	log.Printf("LOGs CREATE ipsec_tunnel value : %v ", d.Get("ipsec_tunnel"))

	if ipSecTunnelPeers, ok := d.GetOk("ipsec_tunnel"); ok {
		clopudIpSecTunnels := ipSecTunnelPeers.(*schema.Set).List()
		// Looping throught List of IPSec Tunnels
		for _, val := range clopudIpSecTunnels {
			ipSecTunnels := val.(map[string]interface{})

			cloudtemplateIpSecTunnelAttr := models.CloudTemplateforIpSectunnelAttributes{}
			cloudtemplateIpSecTunnelAttr.Annotation = "{}"
			cloudtemplateIpSecTunnelAttr.IkeVersion = ipSecTunnels["ike_version"].(string)
			cloudtemplateIpSecTunnelAttr.Poolname = ipSecTunnels["subnet_pool_name"].(string)
			cloudtemplateIpSecTunnelAttr.PreSharedKey = ipSecTunnels["pre_shared_key"].(string)
			cloudtemplateIpSecTunnelAttr.Peeraddr = ipSecTunnels["public_ip_address"].(string)

			cloudtemplateIpSecTunnel := models.NewCloudTemplateforIpSectunnel(fmt.Sprintf(models.RncloudtemplateIpSecTunnel, cloudtemplateIpSecTunnelAttr.Peeraddr), cloudtemplateVpnNetwork.DistinguishedName, cloudtemplateIpSecTunnelAttr)
			err := aciClient.Save(cloudtemplateIpSecTunnel)
			if err != nil {
				return diag.FromErr(err)
			}
			cloudtemplateBgpIpv4Attr := models.BGPIPv4PeerAttributes{}
			cloudtemplateBgpIpv4Attr.Peeraddr = "0.0.0.0/0"
			cloudtemplateBgpIpv4Attr.Peerasn = ipSecTunnels["bgp_peer_asn"].(string)

			cloudtemplateBgpIpv4 := models.NewBGPIPv4Peer(fmt.Sprintf(models.RncloudtemplateBgpIpv4, cloudtemplateBgpIpv4Attr.Peeraddr), cloudtemplateIpSecTunnel.DistinguishedName, cloudtemplateBgpIpv4Attr)
			err = aciClient.Save(cloudtemplateBgpIpv4)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(cloudtemplateVpnNetwork.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudTemplateforVPNNetworkRead(ctx, d, m)
}

func resourceAciCloudTemplateforVPNNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TemplateforVPNNetwork: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TemplateforExternalNetworkDn := d.Get("aci_cloud_external_network_dn").(string)

	cloudtemplateVpnNetworkAttr := models.TemplateforVPNNetworkAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateVpnNetworkAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateVpnNetworkAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudtemplateVpnNetworkAttr.Name = Name.(string)
	}

	if RemoteSiteId, ok := d.GetOk("remote_site_id"); ok {
		cloudtemplateVpnNetworkAttr.RemoteSiteId = RemoteSiteId.(string)
	}

	if RemoteSiteName, ok := d.GetOk("remote_site_name"); ok {
		cloudtemplateVpnNetworkAttr.RemoteSiteName = RemoteSiteName.(string)
	}
	cloudtemplateVpnNetwork := models.NewTemplateforVPNNetwork(fmt.Sprintf(models.RncloudtemplateVpnNetwork, name), TemplateforExternalNetworkDn, nameAlias, cloudtemplateVpnNetworkAttr)

	cloudtemplateVpnNetwork.Status = "modified"

	err := aciClient.Save(cloudtemplateVpnNetwork)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("LOGs UPDATE cloudtemplateVpnNetwork : %v ", cloudtemplateVpnNetwork)
	if ipSecTunnelPeers, ok := d.GetOk("ipsec_tunnel"); ok {
		clopudIpSecTunnels := ipSecTunnelPeers.(*schema.Set).List()
		// Looping throught List of IPSec Tunnels
		for _, val := range clopudIpSecTunnels {
			ipSecTunnels := val.(map[string]interface{})

			cloudtemplateIpSecTunnelAttr := models.CloudTemplateforIpSectunnelAttributes{}
			cloudtemplateIpSecTunnelAttr.Annotation = "{}"
			cloudtemplateIpSecTunnelAttr.IkeVersion = ipSecTunnels["ike_version"].(string)
			cloudtemplateIpSecTunnelAttr.Poolname = ipSecTunnels["subnet_pool_name"].(string)
			cloudtemplateIpSecTunnelAttr.PreSharedKey = ipSecTunnels["pre_shared_key"].(string)
			cloudtemplateIpSecTunnelAttr.Peeraddr = ipSecTunnels["public_ip_address"].(string)

			cloudtemplateIpSecTunnel := models.NewCloudTemplateforIpSectunnel(fmt.Sprintf(models.RncloudtemplateIpSecTunnel, cloudtemplateIpSecTunnelAttr.Peeraddr), cloudtemplateVpnNetwork.DistinguishedName, cloudtemplateIpSecTunnelAttr)
			err := aciClient.Save(cloudtemplateIpSecTunnel)
			if err != nil {
				return diag.FromErr(err)
			}
			cloudtemplateBgpIpv4Attr := models.BGPIPv4PeerAttributes{}
			cloudtemplateBgpIpv4Attr.Peeraddr = "0.0.0.0/0"
			cloudtemplateBgpIpv4Attr.Peerasn = ipSecTunnels["bgp_peer_asn"].(string)

			cloudtemplateBgpIpv4 := models.NewBGPIPv4Peer(fmt.Sprintf(models.RncloudtemplateBgpIpv4, cloudtemplateBgpIpv4Attr.Peeraddr), cloudtemplateIpSecTunnel.DistinguishedName, cloudtemplateBgpIpv4Attr)
			err = aciClient.Save(cloudtemplateBgpIpv4)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(cloudtemplateVpnNetwork.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudTemplateforVPNNetworkRead(ctx, d, m)
}

func resourceAciCloudTemplateforVPNNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	log.Printf("LOGs I ENTER READ ")

	cloudtemplateVpnNetwork, err := getRemoteTemplateforVPNNetwork(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setTemplateforVPNNetworkAttributes(cloudtemplateVpnNetwork, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	cloudtemplateIpSecTunnelData, err := aciClient.ListCloudTemplateforIpSectunnel(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading cloud IPSec Tunnel attributes %v", err)
	}
	cloudtemplateIpSecTunnelSet := make([]map[string]string, 0, 1)
	for _, cloudtemplateIpSecTunnel := range cloudtemplateIpSecTunnelData {

		cloudIpSecTunnelAttMap, cloudtemplateIpSecTunnelDn, err := setCloudTemplateforIpSecTunnelAttributes(cloudtemplateIpSecTunnel, make(map[string]string))
		if err != nil {
			d.SetId("")
			return nil
		}
		log.Printf("LOGs READ cloudIpSecTunnelAttMap 443 : %v ", cloudIpSecTunnelAttMap)

		bgpIPv4PeerData, err := aciClient.ListBGPIPv4Peer(cloudtemplateIpSecTunnelDn)
		if err != nil {
			log.Printf("[DEBUG] Error while reading cloud IPSec Tunnel attributes %v", err)
		}
		log.Printf("LOGs READ FOR bgpIPv4PeerData : %v ", bgpIPv4PeerData)
		for _, bgpIPv4Peer := range bgpIPv4PeerData {
			log.Printf("LOGs READ FOR bgpIPv4Peer : %v ", bgpIPv4Peer)
			bgpPeerAsnAtt, err := getASNfromBGPTPV4Peer(bgpIPv4Peer, make(map[string]string))
			if err != nil {
				d.SetId("")
				return nil
			}
			log.Printf("LOGs READ ASN : %v ", bgpPeerAsnAtt["bgp_peer_asn_att"])
			cloudIpSecTunnelAttMap["bgp_peer_asn"] = bgpPeerAsnAtt["bgp_peer_asn_att"]
			log.Printf("LOGs READ cloudIpSecTunnelAttMap 448 : %v ", cloudIpSecTunnelAttMap)
		}
		log.Printf("LOGs READ cloudIpSecTunnelAttMap 450 : %v ", cloudIpSecTunnelAttMap)
		cloudtemplateIpSecTunnelSet = append(cloudtemplateIpSecTunnelSet, cloudIpSecTunnelAttMap)
		log.Printf("LOGs READ  cloudtemplateIpSecTunnelSet : %v ", cloudtemplateIpSecTunnelSet)
	}
	d.Set("ipsec_tunnel", cloudtemplateIpSecTunnelSet)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudTemplateforVPNNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudtemplateVpnNetwork")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
