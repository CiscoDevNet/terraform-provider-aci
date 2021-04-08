package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBgpPeerConnectivityProfile(addr string, logical_node_profile string, l3_outside string, tenant string, description string, bgpPeerPattr models.BgpPeerConnectivityProfileAttributes) (*models.BgpPeerConnectivityProfile, error) {
	rn := fmt.Sprintf("peerP-[%s]", addr)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", tenant, l3_outside, logical_node_profile)
	bgpPeerP := models.NewBgpPeerConnectivityProfile(rn, parentDn, description, bgpPeerPattr)
	err := sm.Save(bgpPeerP)
	return bgpPeerP, err
}

func (sm *ServiceManager) ReadBgpPeerConnectivityProfile(addr string, logical_node_profile string, l3_outside string, tenant string) (*models.BgpPeerConnectivityProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/peerP-[%s]", tenant, l3_outside, logical_node_profile, addr)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpPeerP := models.BgpPeerConnectivityProfileFromContainer(cont)
	return bgpPeerP, nil
}

func (sm *ServiceManager) DeleteBgpPeerConnectivityProfile(addr string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/peerP-[%s]", tenant, l3_outside, logical_node_profile, addr)
	return sm.DeleteByDn(dn, models.BgppeerpClassName)
}

func (sm *ServiceManager) UpdateBgpPeerConnectivityProfile(addr string, logical_node_profile string, l3_outside string, tenant string, description string, bgpPeerPattr models.BgpPeerConnectivityProfileAttributes) (*models.BgpPeerConnectivityProfile, error) {
	rn := fmt.Sprintf("peerP-[%s]", addr)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", tenant, l3_outside, logical_node_profile)
	bgpPeerP := models.NewBgpPeerConnectivityProfile(rn, parentDn, description, bgpPeerPattr)

	bgpPeerP.Status = "modified"
	err := sm.Save(bgpPeerP)
	return bgpPeerP, err

}

func (sm *ServiceManager) ListBgpPeerConnectivityProfile(logical_node_profile string, l3_outside string, tenant string) ([]*models.BgpPeerConnectivityProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/bgpPeerP.json", baseurlStr, tenant, l3_outside, logical_node_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BgpPeerConnectivityProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationbgpRsPeerPfxPolFromBgpPeerConnectivityProfile(parentDn, tnBgpPeerPfxPolName string) error {
	dn := fmt.Sprintf("%s/rspeerPfxPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnBgpPeerPfxPolName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "bgpRsPeerPfxPol", dn, tnBgpPeerPfxPolName))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cont)

	return nil
}

func (sm *ServiceManager) ReadRelationbgpRsPeerPfxPolFromBgpPeerConnectivityProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "bgpRsPeerPfxPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "bgpRsPeerPfxPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnBgpPeerPfxPolName")
		return dat, err
	} else {
		return nil, err
	}

}
