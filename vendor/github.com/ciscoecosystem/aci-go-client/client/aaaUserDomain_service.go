package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateUserDomain(name string, local_user string, description string, nameAlias string, aaaUserDomainAttr models.UserDomainAttributes) (*models.UserDomain, error) {
	rn := fmt.Sprintf(models.RnaaaUserDomain, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaUserDomain, local_user)
	aaaUserDomain := models.NewUserDomain(rn, parentDn, description, nameAlias, aaaUserDomainAttr)
	err := sm.Save(aaaUserDomain)
	return aaaUserDomain, err
}

func (sm *ServiceManager) ReadUserDomain(name string, local_user string) (*models.UserDomain, error) {
	dn := fmt.Sprintf(models.DnaaaUserDomain, local_user, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaUserDomain := models.UserDomainFromContainer(cont)
	return aaaUserDomain, nil
}

func (sm *ServiceManager) DeleteUserDomain(name string, local_user string) error {
	dn := fmt.Sprintf(models.DnaaaUserDomain, local_user, name)
	return sm.DeleteByDn(dn, models.AaauserdomainClassName)
}

func (sm *ServiceManager) UpdateUserDomain(name string, local_user string, description string, nameAlias string, aaaUserDomainAttr models.UserDomainAttributes) (*models.UserDomain, error) {
	rn := fmt.Sprintf(models.RnaaaUserDomain, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaUserDomain, local_user)
	aaaUserDomain := models.NewUserDomain(rn, parentDn, description, nameAlias, aaaUserDomainAttr)
	aaaUserDomain.Status = "modified"
	err := sm.Save(aaaUserDomain)
	return aaaUserDomain, err
}

func (sm *ServiceManager) ListUserDomain(local_user string) ([]*models.UserDomain, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/user-%s/aaaUserDomain.json", models.BaseurlStr, local_user)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.UserDomainListFromContainer(cont)
	return list, err
}
