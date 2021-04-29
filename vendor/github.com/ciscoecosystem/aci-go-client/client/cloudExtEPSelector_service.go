package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudEndpointSelectorforExternalEPgs(name string, cloud_external_epg string, cloud_application_container string, tenant string, description string, cloudExtEPSelectorattr models.CloudEndpointSelectorforExternalEPgsAttributes) (*models.CloudEndpointSelectorforExternalEPgs, error) {
	rn := fmt.Sprintf("extepselector-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudextepg-%s", tenant, cloud_application_container, cloud_external_epg)
	cloudExtEPSelector := models.NewCloudEndpointSelectorforExternalEPgs(rn, parentDn, description, cloudExtEPSelectorattr)
	err := sm.Save(cloudExtEPSelector)
	return cloudExtEPSelector, err
}

func (sm *ServiceManager) ReadCloudEndpointSelectorforExternalEPgs(name string, cloud_external_epg string, cloud_application_container string, tenant string) (*models.CloudEndpointSelectorforExternalEPgs, error) {
	dn := fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudextepg-%s/extepselector-%s", tenant, cloud_application_container, cloud_external_epg, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudExtEPSelector := models.CloudEndpointSelectorforExternalEPgsFromContainer(cont)
	return cloudExtEPSelector, nil
}

func (sm *ServiceManager) DeleteCloudEndpointSelectorforExternalEPgs(name string, cloud_external_epg string, cloud_application_container string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudextepg-%s/extepselector-%s", tenant, cloud_application_container, cloud_external_epg, name)
	return sm.DeleteByDn(dn, models.CloudextepselectorClassName)
}

func (sm *ServiceManager) UpdateCloudEndpointSelectorforExternalEPgs(name string, cloud_external_epg string, cloud_application_container string, tenant string, description string, cloudExtEPSelectorattr models.CloudEndpointSelectorforExternalEPgsAttributes) (*models.CloudEndpointSelectorforExternalEPgs, error) {
	rn := fmt.Sprintf("extepselector-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudextepg-%s", tenant, cloud_application_container, cloud_external_epg)
	cloudExtEPSelector := models.NewCloudEndpointSelectorforExternalEPgs(rn, parentDn, description, cloudExtEPSelectorattr)

	cloudExtEPSelector.Status = "modified"
	err := sm.Save(cloudExtEPSelector)
	return cloudExtEPSelector, err

}

func (sm *ServiceManager) ListCloudEndpointSelectorforExternalEPgs(cloud_external_epg string, cloud_application_container string, tenant string) ([]*models.CloudEndpointSelectorforExternalEPgs, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/cloudapp-%s/cloudextepg-%s/cloudExtEPSelector.json", baseurlStr, tenant, cloud_application_container, cloud_external_epg)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudEndpointSelectorforExternalEPgsListFromContainer(cont)

	return list, err
}
