package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateL3outRouteTagPolicy(name string, tenant string, description string, l3extRouteTagPolattr models.L3outRouteTagPolicyAttributes) (*models.L3outRouteTagPolicy, error) {
	rn := fmt.Sprintf("rttag-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	l3extRouteTagPol := models.NewL3outRouteTagPolicy(rn, parentDn, description, l3extRouteTagPolattr)
	err := sm.Save(l3extRouteTagPol)
	return l3extRouteTagPol, err
}

func (sm *ServiceManager) ReadL3outRouteTagPolicy(name string, tenant string) (*models.L3outRouteTagPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/rttag-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extRouteTagPol := models.L3outRouteTagPolicyFromContainer(cont)
	return l3extRouteTagPol, nil
}

func (sm *ServiceManager) DeleteL3outRouteTagPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/rttag-%s", tenant, name)
	return sm.DeleteByDn(dn, models.L3extroutetagpolClassName)
}

func (sm *ServiceManager) UpdateL3outRouteTagPolicy(name string, tenant string, description string, l3extRouteTagPolattr models.L3outRouteTagPolicyAttributes) (*models.L3outRouteTagPolicy, error) {
	rn := fmt.Sprintf("rttag-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	l3extRouteTagPol := models.NewL3outRouteTagPolicy(rn, parentDn, description, l3extRouteTagPolattr)

	l3extRouteTagPol.Status = "modified"
	err := sm.Save(l3extRouteTagPol)
	return l3extRouteTagPol, err

}

func (sm *ServiceManager) ListL3outRouteTagPolicy(tenant string) ([]*models.L3outRouteTagPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/l3extRouteTagPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outRouteTagPolicyListFromContainer(cont)

	return list, err
}
