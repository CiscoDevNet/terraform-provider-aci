package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateCloudProviderProfile(vendor string , description string, cloudProvPattr models.CloudProviderProfileAttributes) (*models.CloudProviderProfile, error) {	
	rn := fmt.Sprintf("clouddomp/provp-%s",vendor)
	parentDn := fmt.Sprintf("uni")
	cloudProvP := models.NewCloudProviderProfile(rn, parentDn, description, cloudProvPattr)
	err := sm.Save(cloudProvP)
	return cloudProvP, err
}

func (sm *ServiceManager) ReadCloudProviderProfile(vendor string ) (*models.CloudProviderProfile, error) {
	dn := fmt.Sprintf("uni/clouddomp/provp-%s", vendor )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudProvP := models.CloudProviderProfileFromContainer(cont)
	return cloudProvP, nil
}

func (sm *ServiceManager) DeleteCloudProviderProfile(vendor string ) error {
	dn := fmt.Sprintf("uni/clouddomp/provp-%s", vendor )
	return sm.DeleteByDn(dn, models.CloudprovpClassName)
}

func (sm *ServiceManager) UpdateCloudProviderProfile(vendor string  ,description string, cloudProvPattr models.CloudProviderProfileAttributes) (*models.CloudProviderProfile, error) {
	rn := fmt.Sprintf("clouddomp/provp-%s",vendor)
	parentDn := fmt.Sprintf("uni")
	cloudProvP := models.NewCloudProviderProfile(rn, parentDn, description, cloudProvPattr)

    cloudProvP.Status = "modified"
	err := sm.Save(cloudProvP)
	return cloudProvP, err

}

func (sm *ServiceManager) ListCloudProviderProfile() ([]*models.CloudProviderProfile, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/cloudProvP.json", baseurlStr )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudProviderProfileListFromContainer(cont)

	return list, err
}


