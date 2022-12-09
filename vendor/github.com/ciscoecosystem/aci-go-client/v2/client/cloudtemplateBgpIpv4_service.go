package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateCloudTemplateBGPIPv4Peer(peeraddr string, template_for_ipsec_tunnel_peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string, cloudtemplateBgpIpv4Attr models.CloudTemplateBGPIPv4PeerAttributes) (*models.CloudTemplateBGPIPv4Peer, error) {
	rn := fmt.Sprintf(models.RncloudtemplateBgpIpv4, peeraddr)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateBgpIpv4, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, template_for_ipsec_tunnel_peeraddr)
	cloudtemplateBgpIpv4 := models.NewCloudTemplateBGPIPv4Peer(rn, parentDn, cloudtemplateBgpIpv4Attr)
	err := sm.Save(cloudtemplateBgpIpv4)
	return cloudtemplateBgpIpv4, err
}

func (sm *ServiceManager) ReadCloudTemplateBGPIPv4Peer(peeraddr string, template_for_ipsec_tunnel_peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string) (*models.CloudTemplateBGPIPv4Peer, error) {
	dn := fmt.Sprintf(models.DncloudtemplateBgpIpv4, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, template_for_ipsec_tunnel_peeraddr, peeraddr)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudtemplateBgpIpv4 := models.CloudTemplateBGPIPv4PeerFromContainer(cont)
	return cloudtemplateBgpIpv4, nil
}

func (sm *ServiceManager) DeleteCloudTemplateBGPIPv4Peer(peeraddr string, template_for_ipsec_tunnel_peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string) error {
	dn := fmt.Sprintf(models.DncloudtemplateBgpIpv4, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, template_for_ipsec_tunnel_peeraddr, peeraddr)
	return sm.DeleteByDn(dn, models.Cloudtemplatebgpipv4ClassName)
}

func (sm *ServiceManager) UpdateCloudTemplateBGPIPv4Peer(peeraddr string, template_for_ipsec_tunnel_peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string, cloudtemplateBgpIpv4Attr models.CloudTemplateBGPIPv4PeerAttributes) (*models.CloudTemplateBGPIPv4Peer, error) {
	rn := fmt.Sprintf(models.RncloudtemplateBgpIpv4, peeraddr)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateBgpIpv4, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, template_for_ipsec_tunnel_peeraddr)
	cloudtemplateBgpIpv4 := models.NewCloudTemplateBGPIPv4Peer(rn, parentDn, cloudtemplateBgpIpv4Attr)
	cloudtemplateBgpIpv4.Status = "modified"
	err := sm.Save(cloudtemplateBgpIpv4)
	return cloudtemplateBgpIpv4, err
}

func (sm *ServiceManager) ListCloudTemplateBGPIPv4Peer(parentDn string) ([]*models.CloudTemplateBGPIPv4Peer, error) {
	dnUrl := fmt.Sprintf("%s/%s/cloudtemplateBgpIpv4.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudTemplateBGPIPv4PeerListFromContainer(cont)
	return list, err
}
