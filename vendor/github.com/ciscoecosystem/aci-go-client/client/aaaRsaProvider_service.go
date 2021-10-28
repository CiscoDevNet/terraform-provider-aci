package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRSAProvider(name string, description string, nameAlias string, aaaRsaProviderAttr models.RSAProviderAttributes) (*models.RSAProvider, error) {
	rn := fmt.Sprintf(models.RnaaaRsaProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaRsaProvider)
	aaaRsaProvider := models.NewRSAProvider(rn, parentDn, description, nameAlias, aaaRsaProviderAttr)
	err := sm.Save(aaaRsaProvider)
	return aaaRsaProvider, err
}

func (sm *ServiceManager) ReadRSAProvider(name string) (*models.RSAProvider, error) {
	dn := fmt.Sprintf(models.DnaaaRsaProvider, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaRsaProvider := models.RSAProviderFromContainer(cont)
	return aaaRsaProvider, nil
}

func (sm *ServiceManager) DeleteRSAProvider(name string) error {
	dn := fmt.Sprintf(models.DnaaaRsaProvider, name)
	return sm.DeleteByDn(dn, models.AaarsaproviderClassName)
}

func (sm *ServiceManager) UpdateRSAProvider(name string, description string, nameAlias string, aaaRsaProviderAttr models.RSAProviderAttributes) (*models.RSAProvider, error) {
	rn := fmt.Sprintf(models.RnaaaRsaProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaRsaProvider)
	aaaRsaProvider := models.NewRSAProvider(rn, parentDn, description, nameAlias, aaaRsaProviderAttr)
	aaaRsaProvider.Status = "modified"
	err := sm.Save(aaaRsaProvider)
	return aaaRsaProvider, err
}

func (sm *ServiceManager) ListRSAProvider() ([]*models.RSAProvider, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/rsaext/aaaRsaProvider.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RSAProviderListFromContainer(cont)
	return list, err
}
