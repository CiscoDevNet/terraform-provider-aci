package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudDomainProfile(description string, cloudDomPattr models.CloudDomainProfileAttributes) (*models.CloudDomainProfile, error) {
	rn := fmt.Sprintf("clouddomp")
	parentDn := fmt.Sprintf("uni")
	cloudDomP := models.NewCloudDomainProfile(rn, parentDn, description, cloudDomPattr)
	err := sm.Save(cloudDomP)
	return cloudDomP, err
}

func (sm *ServiceManager) ReadCloudDomainProfile() (*models.CloudDomainProfile, error) {
	dn := fmt.Sprintf("uni/clouddomp")
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudDomP := models.CloudDomainProfileFromContainer(cont)
	return cloudDomP, nil
}

func (sm *ServiceManager) DeleteCloudDomainProfile() error {
	dn := fmt.Sprintf("uni/clouddomp")
	return sm.DeleteByDn(dn, models.ClouddompClassName)
}

func (sm *ServiceManager) UpdateCloudDomainProfile(description string, cloudDomPattr models.CloudDomainProfileAttributes) (*models.CloudDomainProfile, error) {
	rn := fmt.Sprintf("clouddomp")
	parentDn := fmt.Sprintf("uni")
	cloudDomP := models.NewCloudDomainProfile(rn, parentDn, description, cloudDomPattr)

	cloudDomP.Status = "modified"
	err := sm.Save(cloudDomP)
	return cloudDomP, err

}

func (sm *ServiceManager) ListCloudDomainProfile() ([]*models.CloudDomainProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/cloudDomP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudDomainProfileListFromContainer(cont)

	return list, err
}
