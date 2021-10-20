package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTACACSPlusProviderGroup(name string, description string, nameAlias string, aaaTacacsPlusProviderGroupAttr models.TACACSPlusProviderGroupAttributes) (*models.TACACSPlusProviderGroup, error) {
	rn := fmt.Sprintf(models.RnaaaTacacsPlusProviderGroup, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaTacacsPlusProviderGroup)
	aaaTacacsPlusProviderGroup := models.NewTACACSPlusProviderGroup(rn, parentDn, description, nameAlias, aaaTacacsPlusProviderGroupAttr)
	err := sm.Save(aaaTacacsPlusProviderGroup)
	return aaaTacacsPlusProviderGroup, err
}

func (sm *ServiceManager) ReadTACACSPlusProviderGroup(name string) (*models.TACACSPlusProviderGroup, error) {
	dn := fmt.Sprintf(models.DnaaaTacacsPlusProviderGroup, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaTacacsPlusProviderGroup := models.TACACSPlusProviderGroupFromContainer(cont)
	return aaaTacacsPlusProviderGroup, nil
}

func (sm *ServiceManager) DeleteTACACSPlusProviderGroup(name string) error {
	dn := fmt.Sprintf(models.DnaaaTacacsPlusProviderGroup, name)
	return sm.DeleteByDn(dn, models.AaatacacsplusprovidergroupClassName)
}

func (sm *ServiceManager) UpdateTACACSPlusProviderGroup(name string, description string, nameAlias string, aaaTacacsPlusProviderGroupAttr models.TACACSPlusProviderGroupAttributes) (*models.TACACSPlusProviderGroup, error) {
	rn := fmt.Sprintf(models.RnaaaTacacsPlusProviderGroup, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaTacacsPlusProviderGroup)
	aaaTacacsPlusProviderGroup := models.NewTACACSPlusProviderGroup(rn, parentDn, description, nameAlias, aaaTacacsPlusProviderGroupAttr)
	aaaTacacsPlusProviderGroup.Status = "modified"
	err := sm.Save(aaaTacacsPlusProviderGroup)
	return aaaTacacsPlusProviderGroup, err
}

func (sm *ServiceManager) ListTACACSPlusProviderGroup() ([]*models.TACACSPlusProviderGroup, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/tacacsext/aaaTacacsPlusProviderGroup.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.TACACSPlusProviderGroupListFromContainer(cont)
	return list, err
}
