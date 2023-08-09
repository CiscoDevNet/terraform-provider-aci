package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateAciBfdMultihopInterfacePolicy(name string, tenant string, description string, nameAlias string, bfdMhIfPolAttr models.AciBfdMultihopInterfacePolicyAttributes) (*models.AciBfdMultihopInterfacePolicy, error) {
	rn := fmt.Sprintf(models.RnbfdMhIfPol, name)
	parentDn := fmt.Sprintf(models.ParentDnbfdMhIfPol, tenant)
	bfdMhIfPol := models.NewAciBfdMultihopInterfacePolicy(rn, parentDn, description, nameAlias, bfdMhIfPolAttr)
	err := sm.Save(bfdMhIfPol)
	return bfdMhIfPol, err
}

func (sm *ServiceManager) ReadAciBfdMultihopInterfacePolicy(name string, tenant string) (*models.AciBfdMultihopInterfacePolicy, error) {
	dn := fmt.Sprintf(models.DnbfdMhIfPol, tenant, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bfdMhIfPol := models.AciBfdMultihopInterfacePolicyFromContainer(cont)
	return bfdMhIfPol, nil
}

func (sm *ServiceManager) DeleteAciBfdMultihopInterfacePolicy(name string, tenant string) error {
	dn := fmt.Sprintf(models.DnbfdMhIfPol, tenant, name)
	return sm.DeleteByDn(dn, models.BfdmhifpolClassName)
}

func (sm *ServiceManager) UpdateAciBfdMultihopInterfacePolicy(name string, tenant string, description string, nameAlias string, bfdMhIfPolAttr models.AciBfdMultihopInterfacePolicyAttributes) (*models.AciBfdMultihopInterfacePolicy, error) {
	rn := fmt.Sprintf(models.RnbfdMhIfPol, name)
	parentDn := fmt.Sprintf(models.ParentDnbfdMhIfPol, tenant)
	bfdMhIfPol := models.NewAciBfdMultihopInterfacePolicy(rn, parentDn, description, nameAlias, bfdMhIfPolAttr)
	bfdMhIfPol.Status = "modified"
	err := sm.Save(bfdMhIfPol)
	return bfdMhIfPol, err
}

func (sm *ServiceManager) ListAciBfdMultihopInterfacePolicy(tenant string) ([]*models.AciBfdMultihopInterfacePolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/bfdMhIfPol.json", models.BaseurlStr, tenant)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.AciBfdMultihopInterfacePolicyListFromContainer(cont)
	return list, err
}
