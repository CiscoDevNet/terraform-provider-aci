package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateCloudTemplateforVPNNetwork(name string, template_for_external_network string, infra_network_template string, tenant string, nameAlias string, cloudtemplateVpnNetworkAttr models.CloudTemplateforVPNNetworkAttributes) (*models.CloudTemplateforVPNNetwork, error) {
	rn := fmt.Sprintf(models.RncloudtemplateVpnNetwork, name)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateVpnNetwork, tenant, infra_network_template, template_for_external_network)
	cloudtemplateVpnNetwork := models.NewCloudTemplateforVPNNetwork(rn, parentDn, nameAlias, cloudtemplateVpnNetworkAttr)
	err := sm.Save(cloudtemplateVpnNetwork)
	return cloudtemplateVpnNetwork, err
}

func (sm *ServiceManager) ReadCloudTemplateforVPNNetwork(name string, template_for_external_network string, infra_network_template string, tenant string) (*models.CloudTemplateforVPNNetwork, error) {
	dn := fmt.Sprintf(models.DncloudtemplateVpnNetwork, tenant, infra_network_template, template_for_external_network, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudtemplateVpnNetwork := models.CloudTemplateforVPNNetworkFromContainer(cont)
	return cloudtemplateVpnNetwork, nil
}

func (sm *ServiceManager) DeleteCloudTemplateforVPNNetwork(name string, template_for_external_network string, infra_network_template string, tenant string) error {
	dn := fmt.Sprintf(models.DncloudtemplateVpnNetwork, tenant, infra_network_template, template_for_external_network, name)
	return sm.DeleteByDn(dn, models.CloudtemplatevpnnetworkClassName)
}

func (sm *ServiceManager) UpdateCloudTemplateforVPNNetwork(name string, template_for_external_network string, infra_network_template string, tenant string, nameAlias string, cloudtemplateVpnNetworkAttr models.CloudTemplateforVPNNetworkAttributes) (*models.CloudTemplateforVPNNetwork, error) {
	rn := fmt.Sprintf(models.RncloudtemplateVpnNetwork, name)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateVpnNetwork, tenant, infra_network_template, template_for_external_network)
	cloudtemplateVpnNetwork := models.NewCloudTemplateforVPNNetwork(rn, parentDn, nameAlias, cloudtemplateVpnNetworkAttr)
	cloudtemplateVpnNetwork.Status = "modified"
	err := sm.Save(cloudtemplateVpnNetwork)
	return cloudtemplateVpnNetwork, err
}

func (sm *ServiceManager) ListCloudTemplateforVPNNetwork(template_for_external_network string, infra_network_template string, tenant string) ([]*models.CloudTemplateforVPNNetwork, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/infranetwork-%s/extnetwork-%s/cloudtemplateVpnNetwork.json", models.BaseurlStr, tenant, infra_network_template, template_for_external_network)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudTemplateforVPNNetworkListFromContainer(cont)
	return list, err
}
