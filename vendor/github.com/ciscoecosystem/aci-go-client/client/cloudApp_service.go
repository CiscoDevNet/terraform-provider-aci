package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateCloudApplicationcontainer(name string ,tenant string , description string, cloudAppattr models.CloudApplicationcontainerAttributes) (*models.CloudApplicationcontainer, error) {	
	rn := fmt.Sprintf("cloudapp-%s",name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant )
	cloudApp := models.NewCloudApplicationcontainer(rn, parentDn, description, cloudAppattr)
	err := sm.Save(cloudApp)
	return cloudApp, err
}

func (sm *ServiceManager) ReadCloudApplicationcontainer(name string ,tenant string ) (*models.CloudApplicationcontainer, error) {
	dn := fmt.Sprintf("uni/tn-%s/cloudapp-%s", tenant ,name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudApp := models.CloudApplicationcontainerFromContainer(cont)
	return cloudApp, nil
}

func (sm *ServiceManager) DeleteCloudApplicationcontainer(name string ,tenant string ) error {
	dn := fmt.Sprintf("uni/tn-%s/cloudapp-%s", tenant ,name )
	return sm.DeleteByDn(dn, models.CloudappClassName)
}

func (sm *ServiceManager) UpdateCloudApplicationcontainer(name string ,tenant string  ,description string, cloudAppattr models.CloudApplicationcontainerAttributes) (*models.CloudApplicationcontainer, error) {
	rn := fmt.Sprintf("cloudapp-%s",name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant )
	cloudApp := models.NewCloudApplicationcontainer(rn, parentDn, description, cloudAppattr)

    cloudApp.Status = "modified"
	err := sm.Save(cloudApp)
	return cloudApp, err

}

func (sm *ServiceManager) ListCloudApplicationcontainer(tenant string ) ([]*models.CloudApplicationcontainer, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/cloudApp.json", baseurlStr , tenant )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudApplicationcontainerListFromContainer(cont)

	return list, err
}


