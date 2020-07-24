package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateEPGsUsingFunction(tDn string, access_generic string, attachable_access_entity_profile string, description string, infraRsFuncToEpgattr models.EPGsUsingFunctionAttributes) (*models.EPGsUsingFunction, error) {
	rn := fmt.Sprintf("rsfuncToEpg-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/infra/attentp-%s/gen-%s", attachable_access_entity_profile, access_generic)
	infraRsFuncToEpg := models.NewEPGsUsingFunction(rn, parentDn, description, infraRsFuncToEpgattr)
	err := sm.Save(infraRsFuncToEpg)
	return infraRsFuncToEpg, err
}

func (sm *ServiceManager) ReadEPGsUsingFunction(tDn string, access_generic string, attachable_access_entity_profile string) (*models.EPGsUsingFunction, error) {
	dn := fmt.Sprintf("uni/infra/attentp-%s/gen-%s/rsfuncToEpg-[%s]", attachable_access_entity_profile, access_generic, tDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraRsFuncToEpg := models.EPGsUsingFunctionFromContainer(cont)
	return infraRsFuncToEpg, nil
}

func (sm *ServiceManager) DeleteEPGsUsingFunction(tDn string, access_generic string, attachable_access_entity_profile string) error {
	dn := fmt.Sprintf("uni/infra/attentp-%s/gen-%s/rsfuncToEpg-[%s]", attachable_access_entity_profile, access_generic, tDn)
	return sm.DeleteByDn(dn, models.InfrarsfunctoepgClassName)
}

func (sm *ServiceManager) UpdateEPGsUsingFunction(tDn string, access_generic string, attachable_access_entity_profile string, description string, infraRsFuncToEpgattr models.EPGsUsingFunctionAttributes) (*models.EPGsUsingFunction, error) {
	rn := fmt.Sprintf("rsfuncToEpg-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/infra/attentp-%s/gen-%s", attachable_access_entity_profile, access_generic)
	infraRsFuncToEpg := models.NewEPGsUsingFunction(rn, parentDn, description, infraRsFuncToEpgattr)

	infraRsFuncToEpg.Status = "modified"
	err := sm.Save(infraRsFuncToEpg)
	return infraRsFuncToEpg, err

}

func (sm *ServiceManager) ListEPGsUsingFunction(access_generic string, attachable_access_entity_profile string) ([]*models.EPGsUsingFunction, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/attentp-%s/gen-%s/infraRsFuncToEpg.json", baseurlStr, attachable_access_entity_profile, access_generic)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.EPGsUsingFunctionListFromContainer(cont)

	return list, err
}
