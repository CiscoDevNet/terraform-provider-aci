package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func (sm *ServiceManager) CreateVlanEncapsulationforVxlanTraffic(attachable_access_entity_profile string, description string, infraProvAccattr models.VlanEncapsulationforVxlanTrafficAttributes) (*models.VlanEncapsulationforVxlanTraffic, error) {
	rn := fmt.Sprintf("provacc")
	parentDn := fmt.Sprintf("uni/infra/attentp-%s", attachable_access_entity_profile)
	infraProvAcc := models.NewVlanEncapsulationforVxlanTraffic(rn, parentDn, description, infraProvAccattr)
	err := sm.Save(infraProvAcc)
	return infraProvAcc, err
}

func (sm *ServiceManager) ReadVlanEncapsulationforVxlanTraffic(attachable_access_entity_profile string) (*models.VlanEncapsulationforVxlanTraffic, error) {
	dn := fmt.Sprintf("uni/infra/attentp-%s/provacc", attachable_access_entity_profile)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraProvAcc := models.VlanEncapsulationforVxlanTrafficFromContainer(cont)
	return infraProvAcc, nil
}

func (sm *ServiceManager) DeleteVlanEncapsulationforVxlanTraffic(attachable_access_entity_profile string) error {
	dn := fmt.Sprintf("uni/infra/attentp-%s/provacc", attachable_access_entity_profile)
	return sm.DeleteByDn(dn, models.InfraprovaccClassName)
}

func (sm *ServiceManager) UpdateVlanEncapsulationforVxlanTraffic(attachable_access_entity_profile string, description string, infraProvAccattr models.VlanEncapsulationforVxlanTrafficAttributes) (*models.VlanEncapsulationforVxlanTraffic, error) {
	rn := fmt.Sprintf("provacc")
	parentDn := fmt.Sprintf("uni/infra/attentp-%s", attachable_access_entity_profile)
	infraProvAcc := models.NewVlanEncapsulationforVxlanTraffic(rn, parentDn, description, infraProvAccattr)

	infraProvAcc.Status = "modified"
	err := sm.Save(infraProvAcc)
	return infraProvAcc, err

}

func (sm *ServiceManager) ListVlanEncapsulationforVxlanTraffic(attachable_access_entity_profile string) ([]*models.VlanEncapsulationforVxlanTraffic, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/attentp-%s/infraProvAcc.json", baseurlStr, attachable_access_entity_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.VlanEncapsulationforVxlanTrafficListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsFuncToEpgFromVlanEncapsulationforVxlanTraffic(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/provacc/rsfuncToEpg-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "infraRsFuncToEpg", dn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) DeleteRelationinfraRsFuncToEpgFromVlanEncapsulationforVxlanTraffic(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/provacc/rsfuncToEpg-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "infraRsFuncToEpg")
}

func (sm *ServiceManager) ReadRelationinfraRsFuncToEpgFromVlanEncapsulationforVxlanTraffic(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsFuncToEpg")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsFuncToEpg")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
