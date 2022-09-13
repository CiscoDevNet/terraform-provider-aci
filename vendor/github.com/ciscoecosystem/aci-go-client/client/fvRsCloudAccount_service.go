package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTenantToCloudAccountAssociation(tenant string, nameAlias string, fvRsCloudAccountAttr models.TenantToCloudAccountAssociationAttributes) (*models.TenantToCloudAccountAssociation, error) {
	rn := fmt.Sprintf(models.RnfvRsCloudAccount)
	parentDn := fmt.Sprintf(models.ParentDnfvRsCloudAccount, tenant)
	fvRsCloudAccount := models.NewTenantToCloudAccountAssociation(rn, parentDn, nameAlias, fvRsCloudAccountAttr)
	err := sm.Save(fvRsCloudAccount)
	return fvRsCloudAccount, err
}

func (sm *ServiceManager) ReadTenantToCloudAccountAssociation(tenant string) (*models.TenantToCloudAccountAssociation, error) {
	dn := fmt.Sprintf(models.DnfvRsCloudAccount, tenant)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsCloudAccount := models.TenantToCloudAccountAssociationFromContainer(cont)
	return fvRsCloudAccount, nil
}

func (sm *ServiceManager) DeleteTenantToCloudAccountAssociation(tenant string) error {
	dn := fmt.Sprintf(models.DnfvRsCloudAccount, tenant)
	return sm.DeleteByDn(dn, models.FvrscloudaccountClassName)
}

func (sm *ServiceManager) UpdateTenantToCloudAccountAssociation(tenant string, nameAlias string, fvRsCloudAccountAttr models.TenantToCloudAccountAssociationAttributes) (*models.TenantToCloudAccountAssociation, error) {
	rn := fmt.Sprintf(models.RnfvRsCloudAccount)
	parentDn := fmt.Sprintf(models.ParentDnfvRsCloudAccount, tenant)
	fvRsCloudAccount := models.NewTenantToCloudAccountAssociation(rn, parentDn, nameAlias, fvRsCloudAccountAttr)
	fvRsCloudAccount.Status = "modified"
	err := sm.Save(fvRsCloudAccount)
	return fvRsCloudAccount, err
}

func (sm *ServiceManager) ListTenantToCloudAccountAssociation(tenant string) ([]*models.TenantToCloudAccountAssociation, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/fvRsCloudAccount.json", models.BaseurlStr, tenant)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.TenantToCloudAccountAssociationListFromContainer(cont)
	return list, err
}
