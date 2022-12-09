package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateCloudTemplateforIpSectunnel(peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string, cloudtemplateIpSecTunnelAttr models.CloudTemplateforIpSectunnelAttributes) (*models.CloudTemplateforIpSectunnel, error) {
	rn := fmt.Sprintf(models.RncloudtemplateIpSecTunnel, peeraddr)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateIpSecTunnel, tenant, infra_network_template, template_for_external_network, template_for_vpn_network)
	cloudtemplateIpSecTunnel := models.NewCloudTemplateforIpSectunnel(rn, parentDn, cloudtemplateIpSecTunnelAttr)
	err := sm.Save(cloudtemplateIpSecTunnel)
	return cloudtemplateIpSecTunnel, err
}

func (sm *ServiceManager) ReadCloudTemplateforIpSectunnel(peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string) (*models.CloudTemplateforIpSectunnel, error) {
	dn := fmt.Sprintf(models.DncloudtemplateIpSecTunnel, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, peeraddr)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudtemplateIpSecTunnel := models.CloudTemplateforIpSectunnelFromContainer(cont)
	return cloudtemplateIpSecTunnel, nil
}

func (sm *ServiceManager) DeleteCloudTemplateforIpSectunnel(peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string) error {
	dn := fmt.Sprintf(models.DncloudtemplateIpSecTunnel, tenant, infra_network_template, template_for_external_network, template_for_vpn_network, peeraddr)
	return sm.DeleteByDn(dn, models.CloudtemplateipsectunnelClassName)
}

func (sm *ServiceManager) UpdateCloudTemplateforIpSectunnel(peeraddr string, template_for_vpn_network string, template_for_external_network string, infra_network_template string, tenant string, cloudtemplateIpSecTunnelAttr models.CloudTemplateforIpSectunnelAttributes) (*models.CloudTemplateforIpSectunnel, error) {
	rn := fmt.Sprintf(models.RncloudtemplateIpSecTunnel, peeraddr)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateIpSecTunnel, tenant, infra_network_template, template_for_external_network, template_for_vpn_network)
	cloudtemplateIpSecTunnel := models.NewCloudTemplateforIpSectunnel(rn, parentDn, cloudtemplateIpSecTunnelAttr)
	cloudtemplateIpSecTunnel.Status = "modified"
	err := sm.Save(cloudtemplateIpSecTunnel)
	return cloudtemplateIpSecTunnel, err
}

func (sm *ServiceManager) ListCloudTemplateforIpSectunnel(parentDn string) ([]*models.CloudTemplateforIpSectunnel, error) {
	dnUrl := fmt.Sprintf("%s/%s/cloudtemplateIpSecTunnel.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudTemplateforIpSectunnelListFromContainer(cont)
	return list, err
}
