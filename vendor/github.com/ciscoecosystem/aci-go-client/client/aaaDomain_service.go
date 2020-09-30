package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSecurityDomain(name string, description string, aaaDomainattr models.SecurityDomainAttributes) (*models.SecurityDomain, error) {
	rn := fmt.Sprintf("userext/domain-%s", name)
	parentDn := fmt.Sprintf("uni")
	aaaDomain := models.NewSecurityDomain(rn, parentDn, description, aaaDomainattr)
	err := sm.Save(aaaDomain)
	return aaaDomain, err
}

func (sm *ServiceManager) ReadSecurityDomain(name string) (*models.SecurityDomain, error) {
	dn := fmt.Sprintf("uni/userext/domain-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	aaaDomain := models.SecurityDomainFromContainer(cont)
	return aaaDomain, nil
}

func (sm *ServiceManager) DeleteSecurityDomain(name string) error {
	dn := fmt.Sprintf("uni/userext/domain-%s", name)
	return sm.DeleteByDn(dn, models.AaadomainClassName)
}

func (sm *ServiceManager) UpdateSecurityDomain(name string, description string, aaaDomainattr models.SecurityDomainAttributes) (*models.SecurityDomain, error) {
	rn := fmt.Sprintf("userext/domain-%s", name)
	parentDn := fmt.Sprintf("uni")
	aaaDomain := models.NewSecurityDomain(rn, parentDn, description, aaaDomainattr)

	aaaDomain.Status = "modified"
	err := sm.Save(aaaDomain)
	return aaaDomain, err

}

func (sm *ServiceManager) ListSecurityDomain() ([]*models.SecurityDomain, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/aaaDomain.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.SecurityDomainListFromContainer(cont)

	return list, err
}
