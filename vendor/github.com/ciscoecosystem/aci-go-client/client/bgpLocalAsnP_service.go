package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLocalAutonomousSystemProfile(peer_connectivity_profile_addr string, logical_node_profile string, l3_outside string, tenant string, description string, bgpLocalAsnPattr models.LocalAutonomousSystemProfileAttributes) (*models.LocalAutonomousSystemProfile, error) {
	rn := fmt.Sprintf("localasn")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/peerP-[%s]", tenant, l3_outside, logical_node_profile, peer_connectivity_profile_addr)
	bgpLocalAsnP := models.NewLocalAutonomousSystemProfile(rn, parentDn, description, bgpLocalAsnPattr)
	err := sm.Save(bgpLocalAsnP)
	return bgpLocalAsnP, err
}

func (sm *ServiceManager) ReadLocalAutonomousSystemProfile(peer_connectivity_profile_addr string, logical_node_profile string, l3_outside string, tenant string) (*models.LocalAutonomousSystemProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/peerP-[%s]/localasn", tenant, l3_outside, logical_node_profile, peer_connectivity_profile_addr)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpLocalAsnP := models.LocalAutonomousSystemProfileFromContainer(cont)
	return bgpLocalAsnP, nil
}

func (sm *ServiceManager) DeleteLocalAutonomousSystemProfile(peer_connectivity_profile_addr string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/peerP-[%s]/localasn", tenant, l3_outside, logical_node_profile, peer_connectivity_profile_addr)
	return sm.DeleteByDn(dn, models.BgplocalasnpClassName)
}

func (sm *ServiceManager) UpdateLocalAutonomousSystemProfile(peer_connectivity_profile_addr string, logical_node_profile string, l3_outside string, tenant string, description string, bgpLocalAsnPattr models.LocalAutonomousSystemProfileAttributes) (*models.LocalAutonomousSystemProfile, error) {
	rn := fmt.Sprintf("localasn")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/peerP-[%s]", tenant, l3_outside, logical_node_profile, peer_connectivity_profile_addr)
	bgpLocalAsnP := models.NewLocalAutonomousSystemProfile(rn, parentDn, description, bgpLocalAsnPattr)

	bgpLocalAsnP.Status = "modified"
	err := sm.Save(bgpLocalAsnP)
	return bgpLocalAsnP, err

}

func (sm *ServiceManager) ListLocalAutonomousSystemProfile(peer_connectivity_profile_addr string, logical_node_profile string, l3_outside string, tenant string) ([]*models.LocalAutonomousSystemProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/peerP-[%s]/bgpLocalAsnP.json", baseurlStr, tenant, l3_outside, logical_node_profile, peer_connectivity_profile_addr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LocalAutonomousSystemProfileListFromContainer(cont)

	return list, err
}
