package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSAMLProviderGroup(name string, description string, nameAlias string, aaaSamlProviderGroupAttr models.SAMLProviderGroupAttributes) (*models.SAMLProviderGroup, error) {
	rn := fmt.Sprintf(models.RnaaaSamlProviderGroup, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaSamlProviderGroup)
	aaaSamlProviderGroup := models.NewSAMLProviderGroup(rn, parentDn, description, nameAlias, aaaSamlProviderGroupAttr)
	err := sm.Save(aaaSamlProviderGroup)
	return aaaSamlProviderGroup, err
}

func (sm *ServiceManager) ReadSAMLProviderGroup(name string) (*models.SAMLProviderGroup, error) {
	dn := fmt.Sprintf(models.DnaaaSamlProviderGroup, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaSamlProviderGroup := models.SAMLProviderGroupFromContainer(cont)
	return aaaSamlProviderGroup, nil
}

func (sm *ServiceManager) DeleteSAMLProviderGroup(name string) error {
	dn := fmt.Sprintf(models.DnaaaSamlProviderGroup, name)
	return sm.DeleteByDn(dn, models.AaasamlprovidergroupClassName)
}

func (sm *ServiceManager) UpdateSAMLProviderGroup(name string, description string, nameAlias string, aaaSamlProviderGroupAttr models.SAMLProviderGroupAttributes) (*models.SAMLProviderGroup, error) {
	rn := fmt.Sprintf(models.RnaaaSamlProviderGroup, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaSamlProviderGroup)
	aaaSamlProviderGroup := models.NewSAMLProviderGroup(rn, parentDn, description, nameAlias, aaaSamlProviderGroupAttr)
	aaaSamlProviderGroup.Status = "modified"
	err := sm.Save(aaaSamlProviderGroup)
	return aaaSamlProviderGroup, err
}

func (sm *ServiceManager) ListSAMLProviderGroup() ([]*models.SAMLProviderGroup, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/samlext/aaaSamlProviderGroup.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SAMLProviderGroupListFromContainer(cont)
	return list, err
}
