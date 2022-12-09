package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateCloudTemplateforIpSectunnelSourceInterface(sourceInterfaceId string, template_for_ipsec_tunnel_peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string, cloudtemplateIpSecTunnelSourceInterfaceAttr models.CloudTemplateforIpSectunnelSourceInterfaceAttributes) (*models.CloudTemplateforIpSectunnelSourceInterface, error) {
	rn := fmt.Sprintf(models.RncloudtemplateIpSecTunnelSourceInterface, sourceInterfaceId)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateIpSecTunnelSourceInterface, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, template_for_ipsec_tunnel_peeraddr)
	cloudtemplateIpSecTunnelSourceInterface := models.NewCloudTemplateIpSecTunnelSourceInterface(rn, parentDn, cloudtemplateIpSecTunnelSourceInterfaceAttr)
	err := sm.Save(cloudtemplateIpSecTunnelSourceInterface)
	return cloudtemplateIpSecTunnelSourceInterface, err
}

func (sm *ServiceManager) ReadCloudTemplateforIpSectunnelSourceInterface(sourceInterfaceId string, template_for_ipsec_tunnel_peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string) (*models.CloudTemplateforIpSectunnelSourceInterface, error) {
	dn := fmt.Sprintf(models.DncloudtemplateIpSecTunnelSourceInterface, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, template_for_ipsec_tunnel_peeraddr, sourceInterfaceId)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudtemplateIpSecTunnelSourceInterface := models.CloudTemplateIpSecTunnelSourceInterfaceFromContainer(cont)
	return cloudtemplateIpSecTunnelSourceInterface, nil
}

func (sm *ServiceManager) DeleteCloudTemplateforIpSectunnelSourceInterface(sourceInterfaceId string, template_for_ipsec_tunnel_peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string) error {
	dn := fmt.Sprintf(models.DncloudtemplateIpSecTunnelSourceInterface, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, template_for_ipsec_tunnel_peeraddr, sourceInterfaceId)
	return sm.DeleteByDn(dn, models.CloudtemplateipsectunnelsourceinterfaceClassName)
}

func (sm *ServiceManager) UpdateCloudTemplateforIpSectunnelSourceInterface(sourceInterfaceId string, template_for_ipsec_tunnel_peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string, cloudtemplateIpSecTunnelSourceInterfaceAttr models.CloudTemplateforIpSectunnelSourceInterfaceAttributes) (*models.CloudTemplateforIpSectunnelSourceInterface, error) {
	rn := fmt.Sprintf(models.RncloudtemplateIpSecTunnelSourceInterface, sourceInterfaceId)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateIpSecTunnelSourceInterface, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, template_for_ipsec_tunnel_peeraddr)
	cloudtemplateIpSecTunnelSourceInterface := models.NewCloudTemplateIpSecTunnelSourceInterface(rn, parentDn, cloudtemplateIpSecTunnelSourceInterfaceAttr)
	cloudtemplateIpSecTunnelSourceInterface.Status = "modified"
	err := sm.Save(cloudtemplateIpSecTunnelSourceInterface)
	return cloudtemplateIpSecTunnelSourceInterface, err
}

func (sm *ServiceManager) ListCloudTemplateforIpSectunnelSourceInterface(parentDn string) ([]*models.CloudTemplateforIpSectunnelSourceInterface, error) {
	dnUrl := fmt.Sprintf("%s/%s/cloudtemplateIpSecTunnelSourceInterface.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudTemplateIpSecTunnelSourceInterfaceListFromContainer(cont)
	return list, err
}
