package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateHSRPGroupProfile(name string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, hsrpGroupPattr models.HSRPGroupProfileAttributes) (*models.HSRPGroupProfile, error) {
	rn := fmt.Sprintf("hsrpIfP/hsrpGroupP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	hsrpGroupP := models.NewHSRPGroupProfile(rn, parentDn, description, hsrpGroupPattr)
	err := sm.Save(hsrpGroupP)
	return hsrpGroupP, err
}

func (sm *ServiceManager) ReadHSRPGroupProfile(name string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.HSRPGroupProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpGroupP := models.HSRPGroupProfileFromContainer(cont)
	return hsrpGroupP, nil
}

func (sm *ServiceManager) DeleteHSRPGroupProfile(name string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile, name)
	return sm.DeleteByDn(dn, models.HsrpgrouppClassName)
}

func (sm *ServiceManager) UpdateHSRPGroupProfile(name string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, hsrpGroupPattr models.HSRPGroupProfileAttributes) (*models.HSRPGroupProfile, error) {
	rn := fmt.Sprintf("hsrpIfP/hsrpGroupP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	hsrpGroupP := models.NewHSRPGroupProfile(rn, parentDn, description, hsrpGroupPattr)

	hsrpGroupP.Status = "modified"
	err := sm.Save(hsrpGroupP)
	return hsrpGroupP, err

}

func (sm *ServiceManager) ListHSRPGroupProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.HSRPGroupProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpGroupP.json", baseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.HSRPGroupProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationhsrpRsGroupPolFromHSRPGroupProfile(parentDn, tnHsrpGroupPolName string) error {
	dn := fmt.Sprintf("%s/rsGroupPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnHsrpGroupPolName": "%s","annotation":"orchestrator:terraform"
			}
		}
	}`, "hsrpRsGroupPol", dn, tnHsrpGroupPolName))

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

func (sm *ServiceManager) ReadRelationhsrpRsGroupPolFromHSRPGroupProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "hsrpRsGroupPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "hsrpRsGroupPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
