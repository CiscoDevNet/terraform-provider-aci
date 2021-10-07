package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateISISDomainPolicy(name string, description string, nameAlias string, isisDomPolAttr models.ISISDomainPolicyAttributes) (*models.ISISDomainPolicy, error) {
	rn := fmt.Sprintf(models.RnisisDomPol, name)
	parentDn := fmt.Sprintf(models.ParentDnisisDomPol)
	isisDomPol := models.NewISISDomainPolicy(rn, parentDn, description, nameAlias, isisDomPolAttr)
	err := sm.Save(isisDomPol)
	return isisDomPol, err
}

func (sm *ServiceManager) ReadISISDomainPolicy(name string) (*models.ISISDomainPolicy, error) {
	dn := fmt.Sprintf(models.DnisisDomPol, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	isisDomPol := models.ISISDomainPolicyFromContainer(cont)
	return isisDomPol, nil
}

func (sm *ServiceManager) DeleteISISDomainPolicy(name string) error {
	dn := fmt.Sprintf(models.DnisisDomPol, name)
	return sm.DeleteByDn(dn, models.IsisdompolClassName)
}

func (sm *ServiceManager) UpdateISISDomainPolicy(name string, description string, nameAlias string, isisDomPolAttr models.ISISDomainPolicyAttributes) (*models.ISISDomainPolicy, error) {
	rn := fmt.Sprintf(models.RnisisDomPol, name)
	parentDn := fmt.Sprintf(models.ParentDnisisDomPol)
	isisDomPol := models.NewISISDomainPolicy(rn, parentDn, description, nameAlias, isisDomPolAttr)
	isisDomPol.Status = "modified"
	err := sm.Save(isisDomPol)
	return isisDomPol, err
}

func (sm *ServiceManager) ListISISDomainPolicy() ([]*models.ISISDomainPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/isisDomPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ISISDomainPolicyListFromContainer(cont)
	return list, err
}
