package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outBGPProtocolProfile(logical_node_profile string, l3_outside string, tenant string, bgpProtPattr models.L3outBGPProtocolProfileAttributes) (*models.L3outBGPProtocolProfile, error) {
	rn := fmt.Sprintf("protp")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", tenant, l3_outside, logical_node_profile)
	bgpProtP := models.NewL3outBGPProtocolProfile(rn, parentDn, bgpProtPattr)
	err := sm.Save(bgpProtP)
	return bgpProtP, err
}

func (sm *ServiceManager) ReadL3outBGPProtocolProfile(logical_node_profile string, l3_outside string, tenant string) (*models.L3outBGPProtocolProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/protp", tenant, l3_outside, logical_node_profile)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpProtP := models.L3outBGPProtocolProfileFromContainer(cont)
	return bgpProtP, nil
}

func (sm *ServiceManager) DeleteL3outBGPProtocolProfile(logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/protp", tenant, l3_outside, logical_node_profile)
	return sm.DeleteByDn(dn, models.BgpprotpClassName)
}

func (sm *ServiceManager) UpdateL3outBGPProtocolProfile(logical_node_profile string, l3_outside string, tenant string, bgpProtPattr models.L3outBGPProtocolProfileAttributes) (*models.L3outBGPProtocolProfile, error) {
	rn := fmt.Sprintf("protp")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", tenant, l3_outside, logical_node_profile)
	bgpProtP := models.NewL3outBGPProtocolProfile(rn, parentDn, bgpProtPattr)

	bgpProtP.Status = "modified"
	err := sm.Save(bgpProtP)
	return bgpProtP, err

}

func (sm *ServiceManager) ListL3outBGPProtocolProfile(logical_node_profile string, l3_outside string, tenant string) ([]*models.L3outBGPProtocolProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/bgpProtP.json", baseurlStr, tenant, l3_outside, logical_node_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outBGPProtocolProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationbgpRsBgpNodeCtxPolFromL3outBGPProtocolProfile(parentDn, tnBgpCtxPolName string) error {
	dn := fmt.Sprintf("%s/rsbgpNodeCtxPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnBgpCtxPolName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "bgpRsBgpNodeCtxPol", dn, tnBgpCtxPolName))

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

func (sm *ServiceManager) ReadRelationbgpRsBgpNodeCtxPolFromL3outBGPProtocolProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "bgpRsBgpNodeCtxPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "bgpRsBgpNodeCtxPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
