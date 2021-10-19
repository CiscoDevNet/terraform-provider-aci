package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateDefaultAuthenticationMethodforallLogins(description string, nameAlias string, aaaDefaultAuthAttr models.DefaultAuthenticationMethodforallLoginsAttributes) (*models.DefaultAuthenticationMethodforallLogins, error) {
	rn := fmt.Sprintf(models.RnaaaDefaultAuth)
	parentDn := fmt.Sprintf(models.ParentDnaaaDefaultAuth)
	aaaDefaultAuth := models.NewDefaultAuthenticationMethodforallLogins(rn, parentDn, description, nameAlias, aaaDefaultAuthAttr)
	err := sm.Save(aaaDefaultAuth)
	return aaaDefaultAuth, err
}

func (sm *ServiceManager) ReadDefaultAuthenticationMethodforallLogins() (*models.DefaultAuthenticationMethodforallLogins, error) {
	dn := fmt.Sprintf(models.DnaaaDefaultAuth)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaDefaultAuth := models.DefaultAuthenticationMethodforallLoginsFromContainer(cont)
	return aaaDefaultAuth, nil
}

func (sm *ServiceManager) DeleteDefaultAuthenticationMethodforallLogins() error {
	dn := fmt.Sprintf(models.DnaaaDefaultAuth)
	return sm.DeleteByDn(dn, models.AaadefaultauthClassName)
}

func (sm *ServiceManager) UpdateDefaultAuthenticationMethodforallLogins(description string, nameAlias string, aaaDefaultAuthAttr models.DefaultAuthenticationMethodforallLoginsAttributes) (*models.DefaultAuthenticationMethodforallLogins, error) {
	rn := fmt.Sprintf(models.RnaaaDefaultAuth)
	parentDn := fmt.Sprintf(models.ParentDnaaaDefaultAuth)
	aaaDefaultAuth := models.NewDefaultAuthenticationMethodforallLogins(rn, parentDn, description, nameAlias, aaaDefaultAuthAttr)
	aaaDefaultAuth.Status = "modified"
	err := sm.Save(aaaDefaultAuth)
	return aaaDefaultAuth, err
}

func (sm *ServiceManager) ListDefaultAuthenticationMethodforallLogins() ([]*models.DefaultAuthenticationMethodforallLogins, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/authrealm/aaaDefaultAuth.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.DefaultAuthenticationMethodforallLoginsListFromContainer(cont)
	return list, err
}
