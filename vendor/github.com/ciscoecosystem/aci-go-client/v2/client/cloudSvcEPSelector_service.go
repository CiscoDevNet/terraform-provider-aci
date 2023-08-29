package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateCloudServiceEndpointSelector(name string, cloud_service_epg string, cloud_application_container string, tenant string, description string, cloudSvcEPSelectorAttr models.CloudServiceEndpointSelectorAttributes) (*models.CloudServiceEndpointSelector, error) {

	rn := fmt.Sprintf(models.RnCloudSvcEPSelector, name)

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPSelector, tenant, cloud_application_container, cloud_service_epg)
	cloudSvcEPSelector := models.NewCloudServiceEndpointSelector(rn, parentDn, description, cloudSvcEPSelectorAttr)

	err := sm.Save(cloudSvcEPSelector)
	return cloudSvcEPSelector, err
}

func (sm *ServiceManager) ReadCloudServiceEndpointSelector(name string, cloud_service_epg string, cloud_application_container string, tenant string) (*models.CloudServiceEndpointSelector, error) {

	rn := fmt.Sprintf(models.RnCloudSvcEPSelector, name)

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPSelector, tenant, cloud_application_container, cloud_service_epg)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudSvcEPSelector := models.CloudServiceEndpointSelectorFromContainer(cont)
	return cloudSvcEPSelector, nil
}

func (sm *ServiceManager) DeleteCloudServiceEndpointSelector(name string, cloud_service_epg string, cloud_application_container string, tenant string) error {

	rn := fmt.Sprintf(models.RnCloudSvcEPSelector, name)

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPSelector, tenant, cloud_application_container, cloud_service_epg)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	return sm.DeleteByDn(dn, models.CloudSvcEPSelectorClassName)
}

func (sm *ServiceManager) UpdateCloudServiceEndpointSelector(name string, cloud_service_epg string, cloud_application_container string, tenant string, description string, cloudSvcEPSelectorAttr models.CloudServiceEndpointSelectorAttributes) (*models.CloudServiceEndpointSelector, error) {

	rn := fmt.Sprintf(models.RnCloudSvcEPSelector, name)

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPSelector, tenant, cloud_application_container, cloud_service_epg)
	cloudSvcEPSelector := models.NewCloudServiceEndpointSelector(rn, parentDn, description, cloudSvcEPSelectorAttr)

	cloudSvcEPSelector.Status = "modified"
	err := sm.Save(cloudSvcEPSelector)
	return cloudSvcEPSelector, err
}

func (sm *ServiceManager) ListCloudServiceEndpointSelector(cloud_service_epg string, cloud_application_container string, tenant string) ([]*models.CloudServiceEndpointSelector, error) {

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPSelector, tenant, cloud_application_container, cloud_service_epg)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.CloudSvcEPSelectorClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudServiceEndpointSelectorListFromContainer(cont)
	return list, err
}
