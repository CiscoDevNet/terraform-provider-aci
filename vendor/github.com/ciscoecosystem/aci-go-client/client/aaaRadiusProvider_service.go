package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRADIUSProvider(name string, description string, nameAlias string, aaaRadiusProviderAttr models.RADIUSProviderAttributes) (*models.RADIUSProvider, error) {
	rn := fmt.Sprintf(models.RnaaaRadiusProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaRadiusProvider)
	aaaRadiusProvider := models.NewRADIUSProvider(rn, parentDn, description, nameAlias, aaaRadiusProviderAttr)
	err := sm.Save(aaaRadiusProvider)
	return aaaRadiusProvider, err
}

func (sm *ServiceManager) ReadRADIUSProvider(name string) (*models.RADIUSProvider, error) {
	dn := fmt.Sprintf(models.DnaaaRadiusProvider, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaRadiusProvider := models.RADIUSProviderFromContainer(cont)
	return aaaRadiusProvider, nil
}

func (sm *ServiceManager) DeleteRADIUSProvider(name string) error {
	dn := fmt.Sprintf(models.DnaaaRadiusProvider, name)
	return sm.DeleteByDn(dn, models.AaaradiusproviderClassName)
}

func (sm *ServiceManager) UpdateRADIUSProvider(name string, description string, nameAlias string, aaaRadiusProviderAttr models.RADIUSProviderAttributes) (*models.RADIUSProvider, error) {
	rn := fmt.Sprintf(models.RnaaaRadiusProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaRadiusProvider)
	aaaRadiusProvider := models.NewRADIUSProvider(rn, parentDn, description, nameAlias, aaaRadiusProviderAttr)
	aaaRadiusProvider.Status = "modified"
	err := sm.Save(aaaRadiusProvider)
	return aaaRadiusProvider, err
}

func (sm *ServiceManager) ListRADIUSProvider() ([]*models.RADIUSProvider, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/duoext/aaaRadiusProvider.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RADIUSProviderListFromContainer(cont)
	return list, err
}
