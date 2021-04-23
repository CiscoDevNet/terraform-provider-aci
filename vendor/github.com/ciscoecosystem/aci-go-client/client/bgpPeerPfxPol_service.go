package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBGPPeerPrefixPolicy(name string, tenant string, description string, bgpPeerPfxPolattr models.BGPPeerPrefixPolicyAttributes) (*models.BGPPeerPrefixPolicy, error) {
	rn := fmt.Sprintf("bgpPfxP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpPeerPfxPol := models.NewBGPPeerPrefixPolicy(rn, parentDn, description, bgpPeerPfxPolattr)
	err := sm.Save(bgpPeerPfxPol)
	return bgpPeerPfxPol, err
}

func (sm *ServiceManager) ReadBGPPeerPrefixPolicy(name string, tenant string) (*models.BGPPeerPrefixPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/bgpPfxP-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpPeerPfxPol := models.BGPPeerPrefixPolicyFromContainer(cont)
	return bgpPeerPfxPol, nil
}

func (sm *ServiceManager) DeleteBGPPeerPrefixPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/bgpPfxP-%s", tenant, name)
	return sm.DeleteByDn(dn, models.BgppeerpfxpolClassName)
}

func (sm *ServiceManager) UpdateBGPPeerPrefixPolicy(name string, tenant string, description string, bgpPeerPfxPolattr models.BGPPeerPrefixPolicyAttributes) (*models.BGPPeerPrefixPolicy, error) {
	rn := fmt.Sprintf("bgpPfxP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bgpPeerPfxPol := models.NewBGPPeerPrefixPolicy(rn, parentDn, description, bgpPeerPfxPolattr)

	bgpPeerPfxPol.Status = "modified"
	err := sm.Save(bgpPeerPfxPol)
	return bgpPeerPfxPol, err

}

func (sm *ServiceManager) ListBGPPeerPrefixPolicy(tenant string) ([]*models.BGPPeerPrefixPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/bgpPeerPfxPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BGPPeerPrefixPolicyListFromContainer(cont)

	return list, err
}
