package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateDefaultRouteLeakPolicy(parentDn string, l3extDefaultRouteLeakPAttr models.DefaultRouteLeakPolicyAttributes) (*models.DefaultRouteLeakPolicy, error) {
	l3extDefaultRouteLeakP := models.NewDefaultRouteLeakPolicy(models.RnL3extDefaultRouteLeakP, parentDn, l3extDefaultRouteLeakPAttr)
	err := sm.Save(l3extDefaultRouteLeakP)
	return l3extDefaultRouteLeakP, err
}

func (sm *ServiceManager) ReadDefaultRouteLeakPolicy(parentDn string) (*models.DefaultRouteLeakPolicy, error) {
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnL3extDefaultRouteLeakP)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	l3extDefaultRouteLeakP := models.DefaultRouteLeakPolicyFromContainer(cont)
	return l3extDefaultRouteLeakP, nil
}

func (sm *ServiceManager) DeleteDefaultRouteLeakPolicy(parentDn string) error {
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnL3extDefaultRouteLeakP)
	return sm.DeleteByDn(dn, models.L3extDefaultRouteLeakPClassName)
}

func (sm *ServiceManager) UpdateDefaultRouteLeakPolicy(parentDn string, l3extDefaultRouteLeakPAttr models.DefaultRouteLeakPolicyAttributes) (*models.DefaultRouteLeakPolicy, error) {
	l3extDefaultRouteLeakP := models.NewDefaultRouteLeakPolicy(models.RnL3extDefaultRouteLeakP, parentDn, l3extDefaultRouteLeakPAttr)
	l3extDefaultRouteLeakP.Status = "modified"
	err := sm.Save(l3extDefaultRouteLeakP)
	return l3extDefaultRouteLeakP, err
}

func (sm *ServiceManager) ListDefaultRouteLeakPolicy(parentDn string) ([]*models.DefaultRouteLeakPolicy, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.L3extDefaultRouteLeakPClassName)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.DefaultRouteLeakPolicyListFromContainer(cont)
	return list, err
}
