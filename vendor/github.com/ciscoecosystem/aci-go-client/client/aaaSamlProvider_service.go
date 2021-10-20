package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSAMLProvider(name string, description string, nameAlias string, aaaSamlProviderAttr models.SAMLProviderAttributes) (*models.SAMLProvider, error) {
	rn := fmt.Sprintf(models.RnaaaSamlProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaSamlProvider)
	aaaSamlProvider := models.NewSAMLProvider(rn, parentDn, description, nameAlias, aaaSamlProviderAttr)
	err := sm.Save(aaaSamlProvider)
	return aaaSamlProvider, err
}

func (sm *ServiceManager) ReadSAMLProvider(name string) (*models.SAMLProvider, error) {
	dn := fmt.Sprintf(models.DnaaaSamlProvider, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaSamlProvider := models.SAMLProviderFromContainer(cont)
	return aaaSamlProvider, nil
}

func (sm *ServiceManager) DeleteSAMLProvider(name string) error {
	dn := fmt.Sprintf(models.DnaaaSamlProvider, name)
	return sm.DeleteByDn(dn, models.AaasamlproviderClassName)
}

func (sm *ServiceManager) UpdateSAMLProvider(name string, description string, nameAlias string, aaaSamlProviderAttr models.SAMLProviderAttributes) (*models.SAMLProvider, error) {
	rn := fmt.Sprintf(models.RnaaaSamlProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaSamlProvider)
	aaaSamlProvider := models.NewSAMLProvider(rn, parentDn, description, nameAlias, aaaSamlProviderAttr)
	aaaSamlProvider.Status = "modified"
	err := sm.Save(aaaSamlProvider)
	return aaaSamlProvider, err
}

func (sm *ServiceManager) ListSAMLProvider() ([]*models.SAMLProvider, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/samlext/aaaSamlProvider.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SAMLProviderListFromContainer(cont)
	return list, err
}
