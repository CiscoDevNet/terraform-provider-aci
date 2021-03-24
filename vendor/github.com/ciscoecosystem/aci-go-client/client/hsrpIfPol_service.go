package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateHSRPInterfacePolicy(name string, tenant string, description string, hsrpIfPolattr models.HSRPInterfacePolicyAttributes) (*models.HSRPInterfacePolicy, error) {
	rn := fmt.Sprintf("hsrpIfPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	hsrpIfPol := models.NewHSRPInterfacePolicy(rn, parentDn, description, hsrpIfPolattr)
	err := sm.Save(hsrpIfPol)
	return hsrpIfPol, err
}

func (sm *ServiceManager) ReadHSRPInterfacePolicy(name string, tenant string) (*models.HSRPInterfacePolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/hsrpIfPol-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpIfPol := models.HSRPInterfacePolicyFromContainer(cont)
	return hsrpIfPol, nil
}

func (sm *ServiceManager) DeleteHSRPInterfacePolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/hsrpIfPol-%s", tenant, name)
	return sm.DeleteByDn(dn, models.HsrpifpolClassName)
}

func (sm *ServiceManager) UpdateHSRPInterfacePolicy(name string, tenant string, description string, hsrpIfPolattr models.HSRPInterfacePolicyAttributes) (*models.HSRPInterfacePolicy, error) {
	rn := fmt.Sprintf("hsrpIfPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	hsrpIfPol := models.NewHSRPInterfacePolicy(rn, parentDn, description, hsrpIfPolattr)

	hsrpIfPol.Status = "modified"
	err := sm.Save(hsrpIfPol)
	return hsrpIfPol, err

}

func (sm *ServiceManager) ListHSRPInterfacePolicy(tenant string) ([]*models.HSRPInterfacePolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/hsrpIfPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.HSRPInterfacePolicyListFromContainer(cont)

	return list, err
}
