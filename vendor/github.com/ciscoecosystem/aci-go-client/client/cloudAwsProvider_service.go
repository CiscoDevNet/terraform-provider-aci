package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudAWSProvider(tenant string, description string, cloudAwsProviderattr models.CloudAWSProviderAttributes) (*models.CloudAWSProvider, error) {
	rn := fmt.Sprintf("awsprovider")
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	cloudAwsProvider := models.NewCloudAWSProvider(rn, parentDn, description, cloudAwsProviderattr)
	err := sm.Save(cloudAwsProvider)
	return cloudAwsProvider, err
}

func (sm *ServiceManager) ReadCloudAWSProvider(tenant string) (*models.CloudAWSProvider, error) {
	dn := fmt.Sprintf("uni/tn-%s/awsprovider", tenant)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudAwsProvider := models.CloudAWSProviderFromContainer(cont)
	return cloudAwsProvider, nil
}

func (sm *ServiceManager) DeleteCloudAWSProvider(tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/awsprovider", tenant)
	return sm.DeleteByDn(dn, models.CloudawsproviderClassName)
}

func (sm *ServiceManager) UpdateCloudAWSProvider(tenant string, description string, cloudAwsProviderattr models.CloudAWSProviderAttributes) (*models.CloudAWSProvider, error) {
	rn := fmt.Sprintf("awsprovider")
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	cloudAwsProvider := models.NewCloudAWSProvider(rn, parentDn, description, cloudAwsProviderattr)

	cloudAwsProvider.Status = "modified"
	err := sm.Save(cloudAwsProvider)
	return cloudAwsProvider, err

}

func (sm *ServiceManager) ListCloudAWSProvider(tenant string) ([]*models.CloudAWSProvider, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/cloudAwsProvider.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudAWSProviderListFromContainer(cont)

	return list, err
}
