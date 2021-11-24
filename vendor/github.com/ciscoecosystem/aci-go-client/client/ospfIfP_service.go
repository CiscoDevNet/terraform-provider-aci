package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateOSPFInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, ospfIfPattr models.OSPFInterfaceProfileAttributes) (*models.OSPFInterfaceProfile, error) {
	rn := fmt.Sprintf("ospfIfP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	ospfIfP := models.NewOSPFInterfaceProfile(rn, parentDn, description, ospfIfPattr)
	err := sm.Save(ospfIfP)
	return ospfIfP, err
}

func (sm *ServiceManager) ReadOSPFInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.OSPFInterfaceProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/ospfIfP", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfIfP := models.OSPFInterfaceProfileFromContainer(cont)
	return ospfIfP, nil
}

func (sm *ServiceManager) DeleteOSPFInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/ospfIfP", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	return sm.DeleteByDn(dn, models.OspfifpClassName)
}

func (sm *ServiceManager) UpdateOSPFInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, ospfIfPattr models.OSPFInterfaceProfileAttributes) (*models.OSPFInterfaceProfile, error) {
	rn := fmt.Sprintf("ospfIfP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	ospfIfP := models.NewOSPFInterfaceProfile(rn, parentDn, description, ospfIfPattr)

	ospfIfP.Status = "modified"
	err := sm.Save(ospfIfP)
	return ospfIfP, err

}

func (sm *ServiceManager) ListOSPFInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.OSPFInterfaceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/ospfIfP.json", baseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.OSPFInterfaceProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationospfRsIfPolFromInterfaceProfile(parentDn, tnOspfIfPolName string) error {
	dn := fmt.Sprintf("%s/rsIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnOspfIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "ospfRsIfPol", dn, tnOspfIfPolName))

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

	return CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)
}

func (sm *ServiceManager) ReadRelationospfRsIfPolFromInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "ospfRsIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "ospfRsIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
