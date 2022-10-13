package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudTemplateforExternalNetwork(name string, infra_network_template string, tenant string, nameAlias string, cloudtemplateExtNetworkAttr models.CloudTemplateforExternalNetworkAttributes) (*models.CloudTemplateforExternalNetwork, error) {
	rn := fmt.Sprintf(models.RncloudtemplateExtNetwork, name)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateExtNetwork, tenant, infra_network_template)
	cloudtemplateExtNetwork := models.NewCloudTemplateforExternalNetwork(rn, parentDn, nameAlias, cloudtemplateExtNetworkAttr)
	err := sm.Save(cloudtemplateExtNetwork)
	return cloudtemplateExtNetwork, err
}

func (sm *ServiceManager) ReadCloudTemplateforExternalNetwork(name string, infra_network_template string, tenant string) (*models.CloudTemplateforExternalNetwork, error) {
	dn := fmt.Sprintf(models.DncloudtemplateExtNetwork, tenant, infra_network_template, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudtemplateExtNetwork := models.CloudTemplateforExternalNetworkFromContainer(cont)
	return cloudtemplateExtNetwork, nil
}

func (sm *ServiceManager) DeleteCloudTemplateforExternalNetwork(name string, infra_network_template string, tenant string) error {
	dn := fmt.Sprintf(models.DncloudtemplateExtNetwork, tenant, infra_network_template, name)
	return sm.DeleteByDn(dn, models.CloudtemplateextnetworkClassName)
}

func (sm *ServiceManager) UpdateCloudTemplateforExternalNetwork(name string, infra_network_template string, tenant string, nameAlias string, cloudtemplateExtNetworkAttr models.CloudTemplateforExternalNetworkAttributes) (*models.CloudTemplateforExternalNetwork, error) {
	rn := fmt.Sprintf(models.RncloudtemplateExtNetwork, name)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateExtNetwork, tenant, infra_network_template)
	cloudtemplateExtNetwork := models.NewCloudTemplateforExternalNetwork(rn, parentDn, nameAlias, cloudtemplateExtNetworkAttr)
	cloudtemplateExtNetwork.Status = "modified"
	err := sm.Save(cloudtemplateExtNetwork)
	return cloudtemplateExtNetwork, err
}

func (sm *ServiceManager) ListCloudTemplateforExternalNetwork(infra_network_template string, tenant string) ([]*models.CloudTemplateforExternalNetwork, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/infranetwork-%s/cloudtemplateExtNetwork.json", models.BaseurlStr, tenant, infra_network_template)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudTemplateforExternalNetworkListFromContainer(cont)
	return list, err
}
