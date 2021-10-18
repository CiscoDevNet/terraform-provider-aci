package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateAuthenticationMethodfortheDomain(login_domain string, description string, nameAlias string, aaaDomainAuthAttr models.AuthenticationMethodfortheDomainAttributes) (*models.AuthenticationMethodfortheDomain, error) {
	rn := fmt.Sprintf(models.RnaaaDomainAuth)
	parentDn := fmt.Sprintf(models.ParentDnaaaDomainAuth, login_domain)
	aaaDomainAuth := models.NewAuthenticationMethodfortheDomain(rn, parentDn, description, nameAlias, aaaDomainAuthAttr)
	err := sm.Save(aaaDomainAuth)
	return aaaDomainAuth, err
}

func (sm *ServiceManager) ReadAuthenticationMethodfortheDomain(login_domain string) (*models.AuthenticationMethodfortheDomain, error) {
	dn := fmt.Sprintf(models.DnaaaDomainAuth, login_domain)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaDomainAuth := models.AuthenticationMethodfortheDomainFromContainer(cont)
	return aaaDomainAuth, nil
}

func (sm *ServiceManager) DeleteAuthenticationMethodfortheDomain(login_domain string) error {
	dn := fmt.Sprintf(models.DnaaaDomainAuth, login_domain)
	return sm.DeleteByDn(dn, models.AaadomainauthClassName)
}

func (sm *ServiceManager) UpdateAuthenticationMethodfortheDomain(login_domain string, description string, nameAlias string, aaaDomainAuthAttr models.AuthenticationMethodfortheDomainAttributes) (*models.AuthenticationMethodfortheDomain, error) {
	rn := fmt.Sprintf(models.RnaaaDomainAuth)
	parentDn := fmt.Sprintf(models.ParentDnaaaDomainAuth, login_domain)
	aaaDomainAuth := models.NewAuthenticationMethodfortheDomain(rn, parentDn, description, nameAlias, aaaDomainAuthAttr)
	aaaDomainAuth.Status = "modified"
	err := sm.Save(aaaDomainAuth)
	return aaaDomainAuth, err
}

func (sm *ServiceManager) ListAuthenticationMethodfortheDomain(login_domain string) ([]*models.AuthenticationMethodfortheDomain, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/logindomain-%s/aaaDomainAuth.json", models.BaseurlStr, login_domain)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.AuthenticationMethodfortheDomainListFromContainer(cont)
	return list, err
}
