package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateHSRPGroupPolicy(name string, tenant string, description string, hsrpGroupPolattr models.HSRPGroupPolicyAttributes) (*models.HSRPGroupPolicy, error) {
	rn := fmt.Sprintf("hsrpGroupPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	hsrpGroupPol := models.NewHSRPGroupPolicy(rn, parentDn, description, hsrpGroupPolattr)
	err := sm.Save(hsrpGroupPol)
	return hsrpGroupPol, err
}

func (sm *ServiceManager) ReadHSRPGroupPolicy(name string, tenant string) (*models.HSRPGroupPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/hsrpGroupPol-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpGroupPol := models.HSRPGroupPolicyFromContainer(cont)
	return hsrpGroupPol, nil
}

func (sm *ServiceManager) DeleteHSRPGroupPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/hsrpGroupPol-%s", tenant, name)
	return sm.DeleteByDn(dn, models.HsrpgrouppolClassName)
}

func (sm *ServiceManager) UpdateHSRPGroupPolicy(name string, tenant string, description string, hsrpGroupPolattr models.HSRPGroupPolicyAttributes) (*models.HSRPGroupPolicy, error) {
	rn := fmt.Sprintf("hsrpGroupPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	hsrpGroupPol := models.NewHSRPGroupPolicy(rn, parentDn, description, hsrpGroupPolattr)

	hsrpGroupPol.Status = "modified"
	err := sm.Save(hsrpGroupPol)
	return hsrpGroupPol, err

}

func (sm *ServiceManager) ListHSRPGroupPolicy(tenant string) ([]*models.HSRPGroupPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/hsrpGroupPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.HSRPGroupPolicyListFromContainer(cont)

	return list, err
}
