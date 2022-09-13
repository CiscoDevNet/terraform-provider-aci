package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudActiveDirectory(active_directory_id string, tenant string, nameAlias string, cloudADAttr models.CloudActiveDirectoryAttributes) (*models.CloudActiveDirectory, error) {
	rn := fmt.Sprintf(models.RncloudAD, active_directory_id)
	parentDn := fmt.Sprintf(models.ParentDncloudAD, tenant)
	cloudAD := models.NewCloudActiveDirectory(rn, parentDn, nameAlias, cloudADAttr)
	err := sm.Save(cloudAD)
	return cloudAD, err
}

func (sm *ServiceManager) ReadCloudActiveDirectory(active_directory_id string, tenant string) (*models.CloudActiveDirectory, error) {
	dn := fmt.Sprintf(models.DncloudAD, tenant, active_directory_id)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudAD := models.CloudActiveDirectoryFromContainer(cont)
	return cloudAD, nil
}

func (sm *ServiceManager) DeleteCloudActiveDirectory(active_directory_id string, tenant string) error {
	dn := fmt.Sprintf(models.DncloudAD, tenant, active_directory_id)
	return sm.DeleteByDn(dn, models.CloudadClassName)
}

func (sm *ServiceManager) UpdateCloudActiveDirectory(active_directory_id string, tenant string, nameAlias string, cloudADAttr models.CloudActiveDirectoryAttributes) (*models.CloudActiveDirectory, error) {
	rn := fmt.Sprintf(models.RncloudAD, active_directory_id)
	parentDn := fmt.Sprintf(models.ParentDncloudAD, tenant)
	cloudAD := models.NewCloudActiveDirectory(rn, parentDn, nameAlias, cloudADAttr)
	cloudAD.Status = "modified"
	err := sm.Save(cloudAD)
	return cloudAD, err
}

func (sm *ServiceManager) ListCloudActiveDirectory(tenant string) ([]*models.CloudActiveDirectory, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/cloudAD.json", models.BaseurlStr, tenant)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudActiveDirectoryListFromContainer(cont)
	return list, err
}
