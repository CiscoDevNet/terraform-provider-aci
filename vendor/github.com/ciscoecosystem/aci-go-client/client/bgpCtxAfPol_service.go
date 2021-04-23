package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBGPAddressFamilyContextPolicy(name string, tenant string, description string, bgpCtxAfPolattr models.BGPAddressFamilyContextPolicyAttributes) (*models.BGPAddressFamilyContextPolicy, error) {
	rn := fmt.Sprintf("bgpCtxAfP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpCtxAfPol := models.NewBGPAddressFamilyContextPolicy(rn, parentDn, description, bgpCtxAfPolattr)
	err := sm.Save(bgpCtxAfPol)
	return bgpCtxAfPol, err
}

func (sm *ServiceManager) ReadBGPAddressFamilyContextPolicy(name string, tenant string) (*models.BGPAddressFamilyContextPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/bgpCtxAfP-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpCtxAfPol := models.BGPAddressFamilyContextPolicyFromContainer(cont)
	return bgpCtxAfPol, nil
}

func (sm *ServiceManager) DeleteBGPAddressFamilyContextPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/bgpCtxAfP-%s", tenant, name)
	return sm.DeleteByDn(dn, models.BgpctxafpolClassName)
}

func (sm *ServiceManager) UpdateBGPAddressFamilyContextPolicy(name string, tenant string, description string, bgpCtxAfPolattr models.BGPAddressFamilyContextPolicyAttributes) (*models.BGPAddressFamilyContextPolicy, error) {
	rn := fmt.Sprintf("bgpCtxAfP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpCtxAfPol := models.NewBGPAddressFamilyContextPolicy(rn, parentDn, description, bgpCtxAfPolattr)

	bgpCtxAfPol.Status = "modified"
	err := sm.Save(bgpCtxAfPol)
	return bgpCtxAfPol, err

}

func (sm *ServiceManager) ListBGPAddressFamilyContextPolicy(tenant string) ([]*models.BGPAddressFamilyContextPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/bgpCtxAfPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BGPAddressFamilyContextPolicyListFromContainer(cont)

	return list, err
}
