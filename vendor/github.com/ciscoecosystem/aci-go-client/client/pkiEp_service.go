package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreatePublicKeyManagement(description string, nameAlias string, pkiEpAttr models.PublicKeyManagementAttributes) (*models.PublicKeyManagement, error) {
	rn := fmt.Sprintf(models.RnpkiEp)
	parentDn := fmt.Sprintf(models.ParentDnpkiEp)
	pkiEp := models.NewPublicKeyManagement(rn, parentDn, description, nameAlias, pkiEpAttr)
	err := sm.Save(pkiEp)
	return pkiEp, err
}

func (sm *ServiceManager) ReadPublicKeyManagement() (*models.PublicKeyManagement, error) {
	dn := fmt.Sprintf(models.DnpkiEp)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	pkiEp := models.PublicKeyManagementFromContainer(cont)
	return pkiEp, nil
}

func (sm *ServiceManager) DeletePublicKeyManagement() error {
	dn := fmt.Sprintf(models.DnpkiEp)
	return sm.DeleteByDn(dn, models.PkiepClassName)
}

func (sm *ServiceManager) UpdatePublicKeyManagement(description string, nameAlias string, pkiEpAttr models.PublicKeyManagementAttributes) (*models.PublicKeyManagement, error) {
	rn := fmt.Sprintf(models.RnpkiEp)
	parentDn := fmt.Sprintf(models.ParentDnpkiEp)
	pkiEp := models.NewPublicKeyManagement(rn, parentDn, description, nameAlias, pkiEpAttr)
	pkiEp.Status = "modified"
	err := sm.Save(pkiEp)
	return pkiEp, err
}

func (sm *ServiceManager) ListPublicKeyManagement() ([]*models.PublicKeyManagement, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/pkiEp.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.PublicKeyManagementListFromContainer(cont)
	return list, err
}
