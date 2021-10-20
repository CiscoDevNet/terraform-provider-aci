package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRADIUSProviderGroup(name string, description string, nameAlias string, aaaRadiusProviderGroupAttr models.RADIUSProviderGroupAttributes) (*models.RADIUSProviderGroup, error) {
	rn := fmt.Sprintf(models.RnaaaRadiusProviderGroup, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaRadiusProviderGroup)
	aaaRadiusProviderGroup := models.NewRADIUSProviderGroup(rn, parentDn, description, nameAlias, aaaRadiusProviderGroupAttr)
	err := sm.Save(aaaRadiusProviderGroup)
	return aaaRadiusProviderGroup, err
}

func (sm *ServiceManager) ReadRADIUSProviderGroup(name string) (*models.RADIUSProviderGroup, error) {
	dn := fmt.Sprintf(models.DnaaaRadiusProviderGroup, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaRadiusProviderGroup := models.RADIUSProviderGroupFromContainer(cont)
	return aaaRadiusProviderGroup, nil
}

func (sm *ServiceManager) DeleteRADIUSProviderGroup(name string) error {
	dn := fmt.Sprintf(models.DnaaaRadiusProviderGroup, name)
	return sm.DeleteByDn(dn, models.AaaradiusprovidergroupClassName)
}

func (sm *ServiceManager) UpdateRADIUSProviderGroup(name string, description string, nameAlias string, aaaRadiusProviderGroupAttr models.RADIUSProviderGroupAttributes) (*models.RADIUSProviderGroup, error) {
	rn := fmt.Sprintf(models.RnaaaRadiusProviderGroup, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaRadiusProviderGroup)
	aaaRadiusProviderGroup := models.NewRADIUSProviderGroup(rn, parentDn, description, nameAlias, aaaRadiusProviderGroupAttr)
	aaaRadiusProviderGroup.Status = "modified"
	err := sm.Save(aaaRadiusProviderGroup)
	return aaaRadiusProviderGroup, err
}

func (sm *ServiceManager) ListRADIUSProviderGroup() ([]*models.RADIUSProviderGroup, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/radiusext/aaaRadiusProviderGroup.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RADIUSProviderGroupListFromContainer(cont)
	return list, err
}
