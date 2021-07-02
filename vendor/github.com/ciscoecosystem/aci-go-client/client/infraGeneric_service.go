package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateAccessGeneric(name string, attachable_access_entity_profile string, description string, infraGenericattr models.AccessGenericAttributes) (*models.AccessGeneric, error) {
	rn := fmt.Sprintf("gen-%s", name)
	parentDn := fmt.Sprintf("uni/infra/attentp-%s", attachable_access_entity_profile)
	infraGeneric := models.NewAccessGeneric(rn, parentDn, description, infraGenericattr)
	err := sm.Save(infraGeneric)
	return infraGeneric, err
}

func (sm *ServiceManager) ReadAccessGeneric(name string, attachable_access_entity_profile string) (*models.AccessGeneric, error) {
	dn := fmt.Sprintf("uni/infra/attentp-%s/gen-%s", attachable_access_entity_profile, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraGeneric := models.AccessGenericFromContainer(cont)
	return infraGeneric, nil
}

func (sm *ServiceManager) DeleteAccessGeneric(name string, attachable_access_entity_profile string) error {
	dn := fmt.Sprintf("uni/infra/attentp-%s/gen-%s", attachable_access_entity_profile, name)
	return sm.DeleteByDn(dn, models.InfragenericClassName)
}

func (sm *ServiceManager) UpdateAccessGeneric(name string, attachable_access_entity_profile string, description string, infraGenericattr models.AccessGenericAttributes) (*models.AccessGeneric, error) {
	rn := fmt.Sprintf("gen-%s", name)
	parentDn := fmt.Sprintf("uni/infra/attentp-%s", attachable_access_entity_profile)
	infraGeneric := models.NewAccessGeneric(rn, parentDn, description, infraGenericattr)

	infraGeneric.Status = "modified"
	err := sm.Save(infraGeneric)
	return infraGeneric, err

}

func (sm *ServiceManager) ListAccessGeneric(attachable_access_entity_profile string) ([]*models.AccessGeneric, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/attentp-%s/infraGeneric.json", baseurlStr, attachable_access_entity_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.AccessGenericListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsFuncToEpgFromAccessGeneric(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsfuncToEpg-[%s]", parentDn, tDn)
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

	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cont)

	return nil
}

func (sm *ServiceManager) DeleteRelationinfraRsFuncToEpgFromAccessGeneric(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsfuncToEpg-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "infraRsFuncToEpg")
}

func (sm *ServiceManager) ReadRelationinfraRsFuncToEpgFromAccessGeneric(parentDn string) (interface{}, error) {
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
