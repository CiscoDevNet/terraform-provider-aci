package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudEndpointSelector(name string, cloud_epg string, cloud_application_container string, tenant string, description string, cloudEPSelectorattr models.CloudEndpointSelectorAttributes) (*models.CloudEndpointSelector, error) {
	rn := fmt.Sprintf("epselector-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudepg-%s", tenant, cloud_application_container, cloud_epg)
	cloudEPSelector := models.NewCloudEndpointSelector(rn, parentDn, description, cloudEPSelectorattr)
	err := sm.Save(cloudEPSelector)
	return cloudEPSelector, err
}

func (sm *ServiceManager) ReadCloudEndpointSelector(name string, cloud_epg string, cloud_application_container string, tenant string) (*models.CloudEndpointSelector, error) {
	dn := fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudepg-%s/epselector-%s", tenant, cloud_application_container, cloud_epg, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudEPSelector := models.CloudEndpointSelectorFromContainer(cont)
	return cloudEPSelector, nil
}

func (sm *ServiceManager) DeleteCloudEndpointSelector(name string, cloud_epg string, cloud_application_container string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudepg-%s/epselector-%s", tenant, cloud_application_container, cloud_epg, name)
	return sm.DeleteByDn(dn, models.CloudepselectorClassName)
}

func (sm *ServiceManager) UpdateCloudEndpointSelector(name string, cloud_epg string, cloud_application_container string, tenant string, description string, cloudEPSelectorattr models.CloudEndpointSelectorAttributes) (*models.CloudEndpointSelector, error) {
	rn := fmt.Sprintf("epselector-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/cloudapp-%s/cloudepg-%s", tenant, cloud_application_container, cloud_epg)
	cloudEPSelector := models.NewCloudEndpointSelector(rn, parentDn, description, cloudEPSelectorattr)

	cloudEPSelector.Status = "modified"
	err := sm.Save(cloudEPSelector)
	return cloudEPSelector, err

}

func (sm *ServiceManager) ListCloudEndpointSelector(cloud_epg string, cloud_application_container string, tenant string) ([]*models.CloudEndpointSelector, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/cloudapp-%s/cloudepg-%s/cloudEPSelector.json", baseurlStr, tenant, cloud_application_container, cloud_epg)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudEndpointSelectorListFromContainer(cont)

	return list, err
}
