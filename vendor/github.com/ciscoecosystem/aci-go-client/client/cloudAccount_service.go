package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudAccount(vendor string, account_id string, tenant string, nameAlias string, cloudAccountAttr models.CloudAccountAttributes) (*models.CloudAccount, error) {
	rn := fmt.Sprintf(models.RncloudAccount, account_id, vendor)
	parentDn := fmt.Sprintf(models.ParentDncloudAccount, tenant)
	cloudAccount := models.NewCloudAccount(rn, parentDn, nameAlias, cloudAccountAttr)
	err := sm.Save(cloudAccount)
	return cloudAccount, err
}

func (sm *ServiceManager) ReadCloudAccount(vendor string, account_id string, tenant string) (*models.CloudAccount, error) {
	dn := fmt.Sprintf(models.DncloudAccount, tenant, account_id, vendor)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudAccount := models.CloudAccountFromContainer(cont)
	return cloudAccount, nil
}

func (sm *ServiceManager) DeleteCloudAccount(vendor string, account_id string, tenant string) error {
	dn := fmt.Sprintf(models.DncloudAccount, tenant, account_id, vendor)
	return sm.DeleteByDn(dn, models.CloudaccountClassName)
}

func (sm *ServiceManager) UpdateCloudAccount(vendor string, account_id string, tenant string, nameAlias string, cloudAccountAttr models.CloudAccountAttributes) (*models.CloudAccount, error) {
	rn := fmt.Sprintf(models.RncloudAccount, account_id, vendor)
	parentDn := fmt.Sprintf(models.ParentDncloudAccount, tenant)
	cloudAccount := models.NewCloudAccount(rn, parentDn, nameAlias, cloudAccountAttr)
	cloudAccount.Status = "modified"
	err := sm.Save(cloudAccount)
	return cloudAccount, err
}

func (sm *ServiceManager) ListCloudAccount(tenant string) ([]*models.CloudAccount, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/cloudAccount.json", models.BaseurlStr, tenant)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudAccountListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationcloudRsAccountToAccessPolicy(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsaccountToAccessPolicy", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "cloudRsAccountToAccessPolicy", dn, annotation, tDn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationcloudRsAccountToAccessPolicy(parentDn string) error {
	dn := fmt.Sprintf("%s/rsaccountToAccessPolicy", parentDn)
	return sm.DeleteByDn(dn, "cloudRsAccountToAccessPolicy")
}

func (sm *ServiceManager) ReadRelationcloudRsAccountToAccessPolicy(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "cloudRsAccountToAccessPolicy")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "cloudRsAccountToAccessPolicy")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationcloudRsCredentials(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rscredentials", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "cloudRsCredentials", dn, annotation, tDn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationcloudRsCredentials(parentDn string) error {
	dn := fmt.Sprintf("%s/rscredentials", parentDn)
	return sm.DeleteByDn(dn, "cloudRsCredentials")
}

func (sm *ServiceManager) ReadRelationcloudRsCredentials(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "cloudRsCredentials")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "cloudRsCredentials")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
