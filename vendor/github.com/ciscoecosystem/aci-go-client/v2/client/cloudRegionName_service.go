package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateCloudProviderandRegionNames(region string, provider string, template_for_external_network string, infra_network_template string, tenant string, cloudRegionNameAttr models.CloudProviderandRegionNamesAttributes) (*models.CloudProviderandRegionNames, error) {
	rn := fmt.Sprintf(models.RncloudRegionName, provider, region)
	parentDn := fmt.Sprintf(models.ParentDncloudRegionName, tenant, infra_network_template, template_for_external_network)
	cloudRegionName := models.NewCloudProviderandRegionNames(rn, parentDn, cloudRegionNameAttr)
	err := sm.Save(cloudRegionName)
	return cloudRegionName, err
}

func (sm *ServiceManager) ReadCloudProviderandRegionNames(region string, provider string, template_for_external_network string, infra_network_template string, tenant string) (*models.CloudProviderandRegionNames, error) {
	dn := fmt.Sprintf(models.DncloudRegionName, tenant, infra_network_template, template_for_external_network, provider, region)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudRegionName := models.CloudProviderandRegionNamesFromContainer(cont)
	return cloudRegionName, nil
}

func (sm *ServiceManager) DeleteCloudProviderandRegionNames(dn string) error {
	return sm.DeleteByDn(dn, models.CloudregionnameClassName)
}

func (sm *ServiceManager) UpdateCloudProviderandRegionNames(region string, provider string, template_for_external_network string, infra_network_template string, tenant string, cloudRegionNameAttr models.CloudProviderandRegionNamesAttributes) (*models.CloudProviderandRegionNames, error) {
	rn := fmt.Sprintf(models.RncloudRegionName, provider, region)
	parentDn := fmt.Sprintf(models.ParentDncloudRegionName, tenant, infra_network_template, template_for_external_network)
	cloudRegionName := models.NewCloudProviderandRegionNames(rn, parentDn, cloudRegionNameAttr)
	cloudRegionName.Status = "modified"
	err := sm.Save(cloudRegionName)
	return cloudRegionName, err
}

func (sm *ServiceManager) ListCloudProviderandRegionNames(parentDn string) ([]*models.CloudProviderandRegionNames, error) {
	dnUrl := fmt.Sprintf("%s/%s/cloudRegionName.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudProviderandRegionNamesListFromContainer(cont)
	return list, err
}
