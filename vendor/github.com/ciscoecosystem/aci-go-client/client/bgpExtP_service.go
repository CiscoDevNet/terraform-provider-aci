package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outBgpExternalPolicy(l3_outside string, tenant string, description string, bgpExtPattr models.L3outBgpExternalPolicyAttributes) (*models.L3outBgpExternalPolicy, error) {
	rn := fmt.Sprintf("bgpExtP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s", tenant, l3_outside)
	bgpExtP := models.NewL3outBgpExternalPolicy(rn, parentDn, description, bgpExtPattr)
	err := sm.Save(bgpExtP)
	return bgpExtP, err
}

func (sm *ServiceManager) ReadL3outBgpExternalPolicy(l3_outside string, tenant string) (*models.L3outBgpExternalPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/bgpExtP", tenant, l3_outside)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpExtP := models.L3outBgpExternalPolicyFromContainer(cont)
	return bgpExtP, nil
}

func (sm *ServiceManager) DeleteL3outBgpExternalPolicy(l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/bgpExtP", tenant, l3_outside)
	return sm.DeleteByDn(dn, models.BgpextpClassName)
}

func (sm *ServiceManager) UpdateL3outBgpExternalPolicy(l3_outside string, tenant string, description string, bgpExtPattr models.L3outBgpExternalPolicyAttributes) (*models.L3outBgpExternalPolicy, error) {
	rn := fmt.Sprintf("bgpExtP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s", tenant, l3_outside)
	bgpExtP := models.NewL3outBgpExternalPolicy(rn, parentDn, description, bgpExtPattr)

	bgpExtP.Status = "modified"
	err := sm.Save(bgpExtP)
	return bgpExtP, err

}

func (sm *ServiceManager) ListL3outBgpExternalPolicy(l3_outside string, tenant string) ([]*models.L3outBgpExternalPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/bgpExtP.json", baseurlStr, tenant, l3_outside)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outBgpExternalPolicyListFromContainer(cont)

	return list, err
}
