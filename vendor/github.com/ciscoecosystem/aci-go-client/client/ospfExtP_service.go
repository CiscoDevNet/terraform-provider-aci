package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outOspfExternalPolicy(l3_outside string, tenant string, description string, ospfExtPattr models.L3outOspfExternalPolicyAttributes) (*models.L3outOspfExternalPolicy, error) {
	rn := fmt.Sprintf("ospfExtP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s", tenant, l3_outside)
	ospfExtP := models.NewL3outOspfExternalPolicy(rn, parentDn, description, ospfExtPattr)
	err := sm.Save(ospfExtP)
	return ospfExtP, err
}

func (sm *ServiceManager) ReadL3outOspfExternalPolicy(l3_outside string, tenant string) (*models.L3outOspfExternalPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/ospfExtP", tenant, l3_outside)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfExtP := models.L3outOspfExternalPolicyFromContainer(cont)
	return ospfExtP, nil
}

func (sm *ServiceManager) DeleteL3outOspfExternalPolicy(l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/ospfExtP", tenant, l3_outside)
	return sm.DeleteByDn(dn, models.OspfextpClassName)
}

func (sm *ServiceManager) UpdateL3outOspfExternalPolicy(l3_outside string, tenant string, description string, ospfExtPattr models.L3outOspfExternalPolicyAttributes) (*models.L3outOspfExternalPolicy, error) {
	rn := fmt.Sprintf("ospfExtP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s", tenant, l3_outside)
	ospfExtP := models.NewL3outOspfExternalPolicy(rn, parentDn, description, ospfExtPattr)

	ospfExtP.Status = "modified"
	err := sm.Save(ospfExtP)
	return ospfExtP, err

}

func (sm *ServiceManager) ListL3outOspfExternalPolicy(l3_outside string, tenant string) ([]*models.L3outOspfExternalPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/ospfExtP.json", baseurlStr, tenant, l3_outside)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outOspfExternalPolicyListFromContainer(cont)

	return list, err
}
