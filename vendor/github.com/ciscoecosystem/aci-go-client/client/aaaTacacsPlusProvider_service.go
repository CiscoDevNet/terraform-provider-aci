package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTACACSProvider(name string, description string, nameAlias string, aaaTacacsPlusProviderAttr models.TACACSProviderAttributes) (*models.TACACSProvider, error) {
	rn := fmt.Sprintf(models.RnaaaTacacsPlusProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaTacacsPlusProvider)
	aaaTacacsPlusProvider := models.NewTACACSProvider(rn, parentDn, description, nameAlias, aaaTacacsPlusProviderAttr)
	err := sm.Save(aaaTacacsPlusProvider)
	return aaaTacacsPlusProvider, err
}

func (sm *ServiceManager) ReadTACACSProvider(name string) (*models.TACACSProvider, error) {
	dn := fmt.Sprintf(models.DnaaaTacacsPlusProvider, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaTacacsPlusProvider := models.TACACSProviderFromContainer(cont)
	return aaaTacacsPlusProvider, nil
}

func (sm *ServiceManager) DeleteTACACSProvider(name string) error {
	dn := fmt.Sprintf(models.DnaaaTacacsPlusProvider, name)
	return sm.DeleteByDn(dn, models.AaatacacsplusproviderClassName)
}

func (sm *ServiceManager) UpdateTACACSProvider(name string, description string, nameAlias string, aaaTacacsPlusProviderAttr models.TACACSProviderAttributes) (*models.TACACSProvider, error) {
	rn := fmt.Sprintf(models.RnaaaTacacsPlusProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaTacacsPlusProvider)
	aaaTacacsPlusProvider := models.NewTACACSProvider(rn, parentDn, description, nameAlias, aaaTacacsPlusProviderAttr)
	aaaTacacsPlusProvider.Status = "modified"
	err := sm.Save(aaaTacacsPlusProvider)
	return aaaTacacsPlusProvider, err
}

func (sm *ServiceManager) ListTACACSProvider() ([]*models.TACACSProvider, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/tacacsext/aaaTacacsPlusProvider.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.TACACSProviderListFromContainer(cont)
	return list, err
}
