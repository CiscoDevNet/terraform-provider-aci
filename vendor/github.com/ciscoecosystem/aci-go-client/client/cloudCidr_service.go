package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateCloudCIDRPool(addr string ,cloud_context_profile string ,tenant string , description string, cloudCidrattr models.CloudCIDRPoolAttributes) (*models.CloudCIDRPool, error) {	
	rn := fmt.Sprintf("cidr-[%s]",addr)
	parentDn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s", tenant ,cloud_context_profile )
	cloudCidr := models.NewCloudCIDRPool(rn, parentDn, description, cloudCidrattr)
	err := sm.Save(cloudCidr)
	return cloudCidr, err
}

func (sm *ServiceManager) ReadCloudCIDRPool(addr string ,cloud_context_profile string ,tenant string ) (*models.CloudCIDRPool, error) {
	dn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]", tenant ,cloud_context_profile ,addr )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudCidr := models.CloudCIDRPoolFromContainer(cont)
	return cloudCidr, nil
}

func (sm *ServiceManager) DeleteCloudCIDRPool(addr string ,cloud_context_profile string ,tenant string ) error {
	dn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]", tenant ,cloud_context_profile ,addr )
	return sm.DeleteByDn(dn, models.CloudcidrClassName)
}

func (sm *ServiceManager) UpdateCloudCIDRPool(addr string ,cloud_context_profile string ,tenant string  ,description string, cloudCidrattr models.CloudCIDRPoolAttributes) (*models.CloudCIDRPool, error) {
	rn := fmt.Sprintf("cidr-[%s]",addr)
	parentDn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s", tenant ,cloud_context_profile )
	cloudCidr := models.NewCloudCIDRPool(rn, parentDn, description, cloudCidrattr)

    cloudCidr.Status = "modified"
	err := sm.Save(cloudCidr)
	return cloudCidr, err

}

func (sm *ServiceManager) ListCloudCIDRPool(cloud_context_profile string ,tenant string ) ([]*models.CloudCIDRPool, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctxprofile-%s/cloudCidr.json", baseurlStr , tenant ,cloud_context_profile )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudCIDRPoolListFromContainer(cont)

	return list, err
}


