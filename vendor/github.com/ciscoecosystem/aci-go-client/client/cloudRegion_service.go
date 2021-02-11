package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudProvidersRegion(name string, cloud_provider_profile_vendor string, description string, cloudRegionattr models.CloudProvidersRegionAttributes) (*models.CloudProvidersRegion, error) {
	rn := fmt.Sprintf("region-%s", name)
	parentDn := fmt.Sprintf("uni/clouddomp/provp-%s", cloud_provider_profile_vendor)
	cloudRegion := models.NewCloudProvidersRegion(rn, parentDn, description, cloudRegionattr)
	err := sm.Save(cloudRegion)
	return cloudRegion, err
}

func (sm *ServiceManager) ReadCloudProvidersRegion(name string, cloud_provider_profile_vendor string) (*models.CloudProvidersRegion, error) {
	dn := fmt.Sprintf("uni/clouddomp/provp-%s/region-%s", cloud_provider_profile_vendor, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudRegion := models.CloudProvidersRegionFromContainer(cont)
	return cloudRegion, nil
}

func (sm *ServiceManager) DeleteCloudProvidersRegion(name string, cloud_provider_profile_vendor string) error {
	dn := fmt.Sprintf("uni/clouddomp/provp-%s/region-%s", cloud_provider_profile_vendor, name)
	return sm.DeleteByDn(dn, models.CloudregionClassName)
}

func (sm *ServiceManager) UpdateCloudProvidersRegion(name string, cloud_provider_profile_vendor string, description string, cloudRegionattr models.CloudProvidersRegionAttributes) (*models.CloudProvidersRegion, error) {
	rn := fmt.Sprintf("region-%s", name)
	parentDn := fmt.Sprintf("uni/clouddomp/provp-%s", cloud_provider_profile_vendor)
	cloudRegion := models.NewCloudProvidersRegion(rn, parentDn, description, cloudRegionattr)

	cloudRegion.Status = "modified"
	err := sm.Save(cloudRegion)
	return cloudRegion, err

}

func (sm *ServiceManager) ListCloudProvidersRegion(cloud_provider_profile_vendor string) ([]*models.CloudProvidersRegion, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/clouddomp/provp-%s/cloudRegion.json", baseurlStr, cloud_provider_profile_vendor)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudProvidersRegionListFromContainer(cont)

	return list, err
}
