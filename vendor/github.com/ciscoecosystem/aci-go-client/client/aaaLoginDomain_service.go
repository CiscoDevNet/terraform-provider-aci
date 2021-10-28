package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLoginDomain(name string, description string, nameAlias string, aaaLoginDomainAttr models.LoginDomainAttributes) (*models.LoginDomain, error) {
	rn := fmt.Sprintf(models.RnaaaLoginDomain, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLoginDomain)
	aaaLoginDomain := models.NewLoginDomain(rn, parentDn, description, nameAlias, aaaLoginDomainAttr)
	err := sm.Save(aaaLoginDomain)
	return aaaLoginDomain, err
}

func (sm *ServiceManager) ReadLoginDomain(name string) (*models.LoginDomain, error) {
	dn := fmt.Sprintf(models.DnaaaLoginDomain, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLoginDomain := models.LoginDomainFromContainer(cont)
	return aaaLoginDomain, nil
}

func (sm *ServiceManager) DeleteLoginDomain(name string) error {
	dn := fmt.Sprintf(models.DnaaaLoginDomain, name)
	return sm.DeleteByDn(dn, models.AaalogindomainClassName)
}

func (sm *ServiceManager) UpdateLoginDomain(name string, description string, nameAlias string, aaaLoginDomainAttr models.LoginDomainAttributes) (*models.LoginDomain, error) {
	rn := fmt.Sprintf(models.RnaaaLoginDomain, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLoginDomain)
	aaaLoginDomain := models.NewLoginDomain(rn, parentDn, description, nameAlias, aaaLoginDomainAttr)
	aaaLoginDomain.Status = "modified"
	err := sm.Save(aaaLoginDomain)
	return aaaLoginDomain, err
}

func (sm *ServiceManager) ListLoginDomain() ([]*models.LoginDomain, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/aaaLoginDomain.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.LoginDomainListFromContainer(cont)
	return list, err
}
