package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outStaticRouteNextHop(nhAddr string, static_route_ip string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string, description string, ipNexthopPattr models.L3outStaticRouteNextHopAttributes) (*models.L3outStaticRouteNextHop, error) {
	rn := fmt.Sprintf("nh-[%s]", nhAddr)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn, static_route_ip)
	ipNexthopP := models.NewL3outStaticRouteNextHop(rn, parentDn, description, ipNexthopPattr)
	err := sm.Save(ipNexthopP)
	return ipNexthopP, err
}

func (sm *ServiceManager) ReadL3outStaticRouteNextHop(nhAddr string, static_route_ip string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string) (*models.L3outStaticRouteNextHop, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]/nh-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn, static_route_ip, nhAddr)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	ipNexthopP := models.L3outStaticRouteNextHopFromContainer(cont)
	return ipNexthopP, nil
}

func (sm *ServiceManager) DeleteL3outStaticRouteNextHop(nhAddr string, static_route_ip string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]/nh-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn, static_route_ip, nhAddr)
	return sm.DeleteByDn(dn, models.IpnexthoppClassName)
}

func (sm *ServiceManager) UpdateL3outStaticRouteNextHop(nhAddr string, static_route_ip string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string, description string, ipNexthopPattr models.L3outStaticRouteNextHopAttributes) (*models.L3outStaticRouteNextHop, error) {
	rn := fmt.Sprintf("nh-[%s]", nhAddr)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn, static_route_ip)
	ipNexthopP := models.NewL3outStaticRouteNextHop(rn, parentDn, description, ipNexthopPattr)

	ipNexthopP.Status = "modified"
	err := sm.Save(ipNexthopP)
	return ipNexthopP, err

}

func (sm *ServiceManager) ListL3outStaticRouteNextHop(static_route_ip string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string) ([]*models.L3outStaticRouteNextHop, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]/ipNexthopP.json", baseurlStr, tenant, l3_outside, logical_node_profile, fabric_node_tDn, static_route_ip)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outStaticRouteNextHopListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationipRsNexthopRouteTrackFromL3outStaticRouteNextHop(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsNexthopRouteTrack", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "ipRsNexthopRouteTrack", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationipRsNexthopRouteTrackFromL3outStaticRouteNextHop(parentDn string) error {
	dn := fmt.Sprintf("%s/rsNexthopRouteTrack", parentDn)
	return sm.DeleteByDn(dn, "ipRsNexthopRouteTrack")
}

func (sm *ServiceManager) ReadRelationipRsNexthopRouteTrackFromL3outStaticRouteNextHop(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "ipRsNexthopRouteTrack")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "ipRsNexthopRouteTrack")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationipRsNHTrackMemberFromL3outStaticRouteNextHop(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsNHTrackMember", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "ipRsNHTrackMember", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationipRsNHTrackMemberFromL3outStaticRouteNextHop(parentDn string) error {
	dn := fmt.Sprintf("%s/rsNHTrackMember", parentDn)
	return sm.DeleteByDn(dn, "ipRsNHTrackMember")
}

func (sm *ServiceManager) ReadRelationipRsNHTrackMemberFromL3outStaticRouteNextHop(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "ipRsNHTrackMember")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "ipRsNHTrackMember")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
