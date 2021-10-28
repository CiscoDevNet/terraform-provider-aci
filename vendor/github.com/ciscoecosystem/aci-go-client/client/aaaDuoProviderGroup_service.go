package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateDuoProviderGroup(name string, description string, nameAlias string, aaaDuoProviderGroupAttr models.DuoProviderGroupAttributes) (*models.DuoProviderGroup, error) {
	rn := fmt.Sprintf(models.RnaaaDuoProviderGroup, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaDuoProviderGroup)
	aaaDuoProviderGroup := models.NewDuoProviderGroup(rn, parentDn, description, nameAlias, aaaDuoProviderGroupAttr)
	err := sm.Save(aaaDuoProviderGroup)
	return aaaDuoProviderGroup, err
}

func (sm *ServiceManager) ReadDuoProviderGroup(name string) (*models.DuoProviderGroup, error) {
	dn := fmt.Sprintf(models.DnaaaDuoProviderGroup, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaDuoProviderGroup := models.DuoProviderGroupFromContainer(cont)
	return aaaDuoProviderGroup, nil
}

func (sm *ServiceManager) DeleteDuoProviderGroup(name string) error {
	dn := fmt.Sprintf(models.DnaaaDuoProviderGroup, name)
	return sm.DeleteByDn(dn, models.AaaduoprovidergroupClassName)
}

func (sm *ServiceManager) UpdateDuoProviderGroup(name string, description string, nameAlias string, aaaDuoProviderGroupAttr models.DuoProviderGroupAttributes) (*models.DuoProviderGroup, error) {
	rn := fmt.Sprintf(models.RnaaaDuoProviderGroup, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaDuoProviderGroup)
	aaaDuoProviderGroup := models.NewDuoProviderGroup(rn, parentDn, description, nameAlias, aaaDuoProviderGroupAttr)
	aaaDuoProviderGroup.Status = "modified"
	err := sm.Save(aaaDuoProviderGroup)
	return aaaDuoProviderGroup, err
}

func (sm *ServiceManager) ListDuoProviderGroup() ([]*models.DuoProviderGroup, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/duoext/aaaDuoProviderGroup.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.DuoProviderGroupListFromContainer(cont)
	return list, err
}
