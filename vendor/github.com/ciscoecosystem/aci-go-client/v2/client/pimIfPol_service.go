package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreatePIMInterfacePolicy(name string, tenant string, description string, pimIfPolAttr models.PIMInterfacePolicyAttributes) (*models.PIMInterfacePolicy, error) {

	rn := fmt.Sprintf(models.RnPimIfPol, name)

	parentDn := fmt.Sprintf(models.ParentDnPimIfPol, tenant)
	pimIfPol := models.NewPIMInterfacePolicy(rn, parentDn, description, pimIfPolAttr)

	err := sm.Save(pimIfPol)
	return pimIfPol, err
}

func (sm *ServiceManager) ReadPIMInterfacePolicy(name string, tenant string) (*models.PIMInterfacePolicy, error) {

	rn := fmt.Sprintf(models.RnPimIfPol, name)

	parentDn := fmt.Sprintf(models.ParentDnPimIfPol, tenant)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	pimIfPol := models.PIMInterfacePolicyFromContainer(cont)
	return pimIfPol, nil
}

func (sm *ServiceManager) DeletePIMInterfacePolicy(name string, tenant string) error {

	rn := fmt.Sprintf(models.RnPimIfPol, name)

	parentDn := fmt.Sprintf(models.ParentDnPimIfPol, tenant)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	return sm.DeleteByDn(dn, models.PimIfPolClassName)
}

func (sm *ServiceManager) UpdatePIMInterfacePolicy(name string, tenant string, description string, pimIfPolAttr models.PIMInterfacePolicyAttributes) (*models.PIMInterfacePolicy, error) {

	rn := fmt.Sprintf(models.RnPimIfPol, name)

	parentDn := fmt.Sprintf(models.ParentDnPimIfPol, tenant)
	pimIfPol := models.NewPIMInterfacePolicy(rn, parentDn, description, pimIfPolAttr)

	pimIfPol.Status = "modified"
	err := sm.Save(pimIfPol)
	return pimIfPol, err
}

func (sm *ServiceManager) ListPIMInterfacePolicy(tenant string) ([]*models.PIMInterfacePolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimIfPol, tenant)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.PimIfPolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.PIMInterfacePolicyListFromContainer(cont)
	return list, err
}
