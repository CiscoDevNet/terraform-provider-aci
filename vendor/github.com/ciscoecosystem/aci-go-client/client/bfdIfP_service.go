package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBFDInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, bfdIfPattr models.BFDInterfaceProfileAttributes) (*models.BFDInterfaceProfile, error) {
	rn := fmt.Sprintf("bfdIfP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	bfdIfP := models.NewBFDInterfaceProfile(rn, parentDn, description, bfdIfPattr)
	err := sm.Save(bfdIfP)
	return bfdIfP, err
}

func (sm *ServiceManager) ReadBFDInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.BFDInterfaceProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/bfdIfP", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bfdIfP := models.BFDInterfaceProfileFromContainer(cont)
	return bfdIfP, nil
}

func (sm *ServiceManager) DeleteBFDInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/bfdIfP", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	return sm.DeleteByDn(dn, models.BfdifpClassName)
}

func (sm *ServiceManager) UpdateBFDInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, bfdIfPattr models.BFDInterfaceProfileAttributes) (*models.BFDInterfaceProfile, error) {
	rn := fmt.Sprintf("bfdIfP")
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	bfdIfP := models.NewBFDInterfaceProfile(rn, parentDn, description, bfdIfPattr)

	bfdIfP.Status = "modified"
	err := sm.Save(bfdIfP)
	return bfdIfP, err

}

func (sm *ServiceManager) ListBFDInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.BFDInterfaceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/bfdIfP.json", baseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BFDInterfaceProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationbfdRsIfPolFromInterfaceProfile(parentDn, tnBfdIfPolName string) error {
	dn := fmt.Sprintf("%s/bfdIfP/rsIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnBfdIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "bfdRsIfPol", dn, tnBfdIfPolName))

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

func (sm *ServiceManager) ReadRelationbfdRsIfPolFromInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "bfdRsIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "bfdRsIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
