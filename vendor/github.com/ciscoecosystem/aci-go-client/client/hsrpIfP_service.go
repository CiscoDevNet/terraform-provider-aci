package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outHSRPInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, hsrpIfPattr models.L3outHSRPInterfaceProfileAttributes) (*models.L3outHSRPInterfaceProfile, error) {
	rn := fmt.Sprintf("hsrpIfP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	hsrpIfP := models.NewL3outHSRPInterfaceProfile(rn, parentDn, description, hsrpIfPattr)
	err := sm.Save(hsrpIfP)
	return hsrpIfP, err
}

func (sm *ServiceManager) ReadL3outHSRPInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.L3outHSRPInterfaceProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpIfP := models.L3outHSRPInterfaceProfileFromContainer(cont)
	return hsrpIfP, nil
}

func (sm *ServiceManager) DeleteL3outHSRPInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	return sm.DeleteByDn(dn, models.HsrpifpClassName)
}

func (sm *ServiceManager) UpdateL3outHSRPInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, hsrpIfPattr models.L3outHSRPInterfaceProfileAttributes) (*models.L3outHSRPInterfaceProfile, error) {
	rn := fmt.Sprintf("hsrpIfP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	hsrpIfP := models.NewL3outHSRPInterfaceProfile(rn, parentDn, description, hsrpIfPattr)

	hsrpIfP.Status = "modified"
	err := sm.Save(hsrpIfP)
	return hsrpIfP, err

}

func (sm *ServiceManager) ListL3outHSRPInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.L3outHSRPInterfaceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP.json", baseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outHSRPInterfaceProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationhsrpRsIfPolFromL3outHSRPInterfaceProfile(parentDn, tnHsrpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnHsrpIfPolName": "%s", 
				"annotation":"orchestrator:terraform"
			}
		}
	}`, "hsrpRsIfPol", dn, tnHsrpIfPolName))

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

func (sm *ServiceManager) ReadRelationhsrpRsIfPolFromL3outHSRPInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "hsrpRsIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "hsrpRsIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
