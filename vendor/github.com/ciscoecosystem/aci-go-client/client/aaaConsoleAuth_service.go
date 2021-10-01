package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateConsoleAuthenticationMethod(description string, nameAlias string, aaaConsoleAuthAttr models.ConsoleAuthenticationMethodAttributes) (*models.ConsoleAuthenticationMethod, error) {
	rn := fmt.Sprintf(models.RnaaaConsoleAuth)
	parentDn := fmt.Sprintf(models.ParentDnaaaConsoleAuth)
	aaaConsoleAuth := models.NewConsoleAuthenticationMethod(rn, parentDn, description, nameAlias, aaaConsoleAuthAttr)
	err := sm.Save(aaaConsoleAuth)
	return aaaConsoleAuth, err
}

func (sm *ServiceManager) ReadConsoleAuthenticationMethod() (*models.ConsoleAuthenticationMethod, error) {
	dn := fmt.Sprintf(models.DnaaaConsoleAuth)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaConsoleAuth := models.ConsoleAuthenticationMethodFromContainer(cont)
	return aaaConsoleAuth, nil
}

func (sm *ServiceManager) DeleteConsoleAuthenticationMethod() error {
	dn := fmt.Sprintf(models.DnaaaConsoleAuth)
	return sm.DeleteByDn(dn, models.AaaconsoleauthClassName)
}

func (sm *ServiceManager) UpdateConsoleAuthenticationMethod(description string, nameAlias string, aaaConsoleAuthAttr models.ConsoleAuthenticationMethodAttributes) (*models.ConsoleAuthenticationMethod, error) {
	rn := fmt.Sprintf(models.RnaaaConsoleAuth)
	parentDn := fmt.Sprintf(models.ParentDnaaaConsoleAuth)
	aaaConsoleAuth := models.NewConsoleAuthenticationMethod(rn, parentDn, description, nameAlias, aaaConsoleAuthAttr)
	aaaConsoleAuth.Status = "modified"
	err := sm.Save(aaaConsoleAuth)
	return aaaConsoleAuth, err
}

func (sm *ServiceManager) ListConsoleAuthenticationMethod() ([]*models.ConsoleAuthenticationMethod, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/authrealm/aaaConsoleAuth.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ConsoleAuthenticationMethodListFromContainer(cont)
	return list, err
}
