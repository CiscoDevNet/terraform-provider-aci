package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudAvailabilityZone(name string, cloud_providers_region string, cloud_provider_profile_vendor string, description string, cloudZoneattr models.CloudAvailabilityZoneAttributes) (*models.CloudAvailabilityZone, error) {
	rn := fmt.Sprintf("zone-%s", name)
	parentDn := fmt.Sprintf("uni/clouddomp/provp-%s/region-%s", cloud_provider_profile_vendor, cloud_providers_region)
	cloudZone := models.NewCloudAvailabilityZone(rn, parentDn, description, cloudZoneattr)
	err := sm.Save(cloudZone)
	return cloudZone, err
}

func (sm *ServiceManager) ReadCloudAvailabilityZone(name string, cloud_providers_region string, cloud_provider_profile_vendor string) (*models.CloudAvailabilityZone, error) {
	dn := fmt.Sprintf("uni/clouddomp/provp-%s/region-%s/zone-%s", cloud_provider_profile_vendor, cloud_providers_region, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudZone := models.CloudAvailabilityZoneFromContainer(cont)
	return cloudZone, nil
}

func (sm *ServiceManager) DeleteCloudAvailabilityZone(name string, cloud_providers_region string, cloud_provider_profile_vendor string) error {
	dn := fmt.Sprintf("uni/clouddomp/provp-%s/region-%s/zone-%s", cloud_provider_profile_vendor, cloud_providers_region, name)
	return sm.DeleteByDn(dn, models.CloudzoneClassName)
}

func (sm *ServiceManager) UpdateCloudAvailabilityZone(name string, cloud_providers_region string, cloud_provider_profile_vendor string, description string, cloudZoneattr models.CloudAvailabilityZoneAttributes) (*models.CloudAvailabilityZone, error) {
	rn := fmt.Sprintf("zone-%s", name)
	parentDn := fmt.Sprintf("uni/clouddomp/provp-%s/region-%s", cloud_provider_profile_vendor, cloud_providers_region)
	cloudZone := models.NewCloudAvailabilityZone(rn, parentDn, description, cloudZoneattr)

	cloudZone.Status = "modified"
	err := sm.Save(cloudZone)
	return cloudZone, err

}

func (sm *ServiceManager) ListCloudAvailabilityZone(cloud_providers_region string, cloud_provider_profile_vendor string) ([]*models.CloudAvailabilityZone, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/clouddomp/provp-%s/region-%s/cloudZone.json", baseurlStr, cloud_provider_profile_vendor, cloud_providers_region)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudAvailabilityZoneListFromContainer(cont)

	return list, err
}
