package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateOSPFTimersPolicy(name string, tenant string, description string, ospfCtxPolattr models.OSPFTimersPolicyAttributes) (*models.OSPFTimersPolicy, error) {
	rn := fmt.Sprintf("ospfCtxP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	ospfCtxPol := models.NewOSPFTimersPolicy(rn, parentDn, description, ospfCtxPolattr)
	err := sm.Save(ospfCtxPol)
	return ospfCtxPol, err
}

func (sm *ServiceManager) ReadOSPFTimersPolicy(name string, tenant string) (*models.OSPFTimersPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/ospfCtxP-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfCtxPol := models.OSPFTimersPolicyFromContainer(cont)
	return ospfCtxPol, nil
}

func (sm *ServiceManager) DeleteOSPFTimersPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ospfCtxP-%s", tenant, name)
	return sm.DeleteByDn(dn, models.OspfctxpolClassName)
}

func (sm *ServiceManager) UpdateOSPFTimersPolicy(name string, tenant string, description string, ospfCtxPolattr models.OSPFTimersPolicyAttributes) (*models.OSPFTimersPolicy, error) {
	rn := fmt.Sprintf("ospfCtxP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	ospfCtxPol := models.NewOSPFTimersPolicy(rn, parentDn, description, ospfCtxPolattr)

	ospfCtxPol.Status = "modified"
	err := sm.Save(ospfCtxPol)
	return ospfCtxPol, err

}

func (sm *ServiceManager) ListOSPFTimersPolicy(tenant string) ([]*models.OSPFTimersPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ospfCtxPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.OSPFTimersPolicyListFromContainer(cont)

	return list, err
}
