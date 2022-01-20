package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateWebTokenData(description string, nameAlias string, pkiWebTokenDataAttr models.WebTokenDataAttributes) (*models.WebTokenData, error) {
	rn := fmt.Sprintf(models.RnpkiWebTokenData)
	parentDn := fmt.Sprintf(models.ParentDnpkiWebTokenData)
	pkiWebTokenData := models.NewWebTokenData(rn, parentDn, description, nameAlias, pkiWebTokenDataAttr)
	err := sm.Save(pkiWebTokenData)
	return pkiWebTokenData, err
}

func (sm *ServiceManager) ReadWebTokenData() (*models.WebTokenData, error) {
	dn := fmt.Sprintf(models.DnpkiWebTokenData)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	pkiWebTokenData := models.WebTokenDataFromContainer(cont)
	return pkiWebTokenData, nil
}

func (sm *ServiceManager) DeleteWebTokenData() error {
	dn := fmt.Sprintf(models.DnpkiWebTokenData)
	return sm.DeleteByDn(dn, models.PkiwebtokendataClassName)
}

func (sm *ServiceManager) UpdateWebTokenData(description string, nameAlias string, pkiWebTokenDataAttr models.WebTokenDataAttributes) (*models.WebTokenData, error) {
	rn := fmt.Sprintf(models.RnpkiWebTokenData)
	parentDn := fmt.Sprintf(models.ParentDnpkiWebTokenData)
	pkiWebTokenData := models.NewWebTokenData(rn, parentDn, description, nameAlias, pkiWebTokenDataAttr)
	pkiWebTokenData.Status = "modified"
	err := sm.Save(pkiWebTokenData)
	return pkiWebTokenData, err
}

func (sm *ServiceManager) ListWebTokenData() ([]*models.WebTokenData, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/pkiext/pkiWebTokenData.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.WebTokenDataListFromContainer(cont)
	return list, err
}
