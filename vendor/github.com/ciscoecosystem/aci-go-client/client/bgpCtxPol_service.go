package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBGPTimersPolicy(name string, tenant string, description string, bgpCtxPolattr models.BGPTimersPolicyAttributes) (*models.BGPTimersPolicy, error) {
	rn := fmt.Sprintf("bgpCtxP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpCtxPol := models.NewBGPTimersPolicy(rn, parentDn, description, bgpCtxPolattr)
	err := sm.Save(bgpCtxPol)
	return bgpCtxPol, err
}

func (sm *ServiceManager) ReadBGPTimersPolicy(name string, tenant string) (*models.BGPTimersPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/bgpCtxP-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpCtxPol := models.BGPTimersPolicyFromContainer(cont)
	return bgpCtxPol, nil
}

func (sm *ServiceManager) DeleteBGPTimersPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/bgpCtxP-%s", tenant, name)
	return sm.DeleteByDn(dn, models.BgpctxpolClassName)
}

func (sm *ServiceManager) UpdateBGPTimersPolicy(name string, tenant string, description string, bgpCtxPolattr models.BGPTimersPolicyAttributes) (*models.BGPTimersPolicy, error) {
	rn := fmt.Sprintf("bgpCtxP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpCtxPol := models.NewBGPTimersPolicy(rn, parentDn, description, bgpCtxPolattr)

	bgpCtxPol.Status = "modified"
	err := sm.Save(bgpCtxPol)
	return bgpCtxPol, err

}

func (sm *ServiceManager) ListBGPTimersPolicy(tenant string) ([]*models.BGPTimersPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/bgpCtxPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BGPTimersPolicyListFromContainer(cont)

	return list, err
}
