package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outStaticRoute(ip string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string, description string, ipRoutePattr models.L3outStaticRouteAttributes) (*models.L3outStaticRoute, error) {
	rn := fmt.Sprintf("rt-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn)
	ipRouteP := models.NewL3outStaticRoute(rn, parentDn, description, ipRoutePattr)
	err := sm.Save(ipRouteP)
	return ipRouteP, err
}

func (sm *ServiceManager) ReadL3outStaticRoute(ip string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string) (*models.L3outStaticRoute, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn, ip)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	ipRouteP := models.L3outStaticRouteFromContainer(cont)
	return ipRouteP, nil
}

func (sm *ServiceManager) DeleteL3outStaticRoute(ip string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn, ip)
	return sm.DeleteByDn(dn, models.IproutepClassName)
}

func (sm *ServiceManager) UpdateL3outStaticRoute(ip string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string, description string, ipRoutePattr models.L3outStaticRouteAttributes) (*models.L3outStaticRoute, error) {
	rn := fmt.Sprintf("rt-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn)
	ipRouteP := models.NewL3outStaticRoute(rn, parentDn, description, ipRoutePattr)

	ipRouteP.Status = "modified"
	err := sm.Save(ipRouteP)
	return ipRouteP, err

}

func (sm *ServiceManager) ListL3outStaticRoute(fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string) ([]*models.L3outStaticRoute, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/ipRouteP.json", baseurlStr, tenant, l3_outside, logical_node_profile, fabric_node_tDn)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outStaticRouteListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationipRsRouteTrackFromL3outStaticRoute(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsRouteTrack", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "ipRsRouteTrack", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationipRsRouteTrackFromL3outStaticRoute(parentDn string) error {
	dn := fmt.Sprintf("%s/rsRouteTrack", parentDn)
	return sm.DeleteByDn(dn, "ipRsRouteTrack")
}

func (sm *ServiceManager) ReadRelationipRsRouteTrackFromL3outStaticRoute(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "ipRsRouteTrack")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "ipRsRouteTrack")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
