package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateIGMPInterfacePolicy(name string, tenant string, description string, igmpIfPolAttr models.IGMPInterfacePolicyAttributes) (*models.IGMPInterfacePolicy, error) {

	rn := fmt.Sprintf(models.RnIgmpIfPol, name)

	parentDn := fmt.Sprintf(models.ParentDnIgmpIfPol, tenant)
	igmpIfPol := models.NewIGMPInterfacePolicy(rn, parentDn, description, igmpIfPolAttr)

	err := sm.Save(igmpIfPol)
	return igmpIfPol, err
}

func (sm *ServiceManager) ReadIGMPInterfacePolicy(name string, tenant string) (*models.IGMPInterfacePolicy, error) {

	rn := fmt.Sprintf(models.RnIgmpIfPol, name)

	parentDn := fmt.Sprintf(models.ParentDnIgmpIfPol, tenant)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	igmpIfPol := models.IGMPInterfacePolicyFromContainer(cont)
	return igmpIfPol, nil
}

func (sm *ServiceManager) DeleteIGMPInterfacePolicy(name string, tenant string) error {

	rn := fmt.Sprintf(models.RnIgmpIfPol, name)

	parentDn := fmt.Sprintf(models.ParentDnIgmpIfPol, tenant)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	return sm.DeleteByDn(dn, models.IgmpIfPolClassName)
}

func (sm *ServiceManager) UpdateIGMPInterfacePolicy(name string, tenant string, description string, igmpIfPolAttr models.IGMPInterfacePolicyAttributes) (*models.IGMPInterfacePolicy, error) {

	rn := fmt.Sprintf(models.RnIgmpIfPol, name)

	parentDn := fmt.Sprintf(models.ParentDnIgmpIfPol, tenant)
	igmpIfPol := models.NewIGMPInterfacePolicy(rn, parentDn, description, igmpIfPolAttr)

	igmpIfPol.Status = "modified"
	err := sm.Save(igmpIfPol)
	return igmpIfPol, err
}

func (sm *ServiceManager) ListIGMPInterfacePolicy(tenant string) ([]*models.IGMPInterfacePolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpIfPol, tenant)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.IgmpIfPolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.IGMPInterfacePolicyListFromContainer(cont)
	return list, err
}
