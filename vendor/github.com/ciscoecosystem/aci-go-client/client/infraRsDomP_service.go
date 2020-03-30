package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateDomain(tDn string, attachable_access_entity_profile string, description string, infraRsDomPattr models.DomainAttributes) (*models.Domain, error) {
	rn := fmt.Sprintf("rsdomP-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/infra/attentp-%s", attachable_access_entity_profile)
	infraRsDomP := models.NewDomain(rn, parentDn, description, infraRsDomPattr)
	err := sm.Save(infraRsDomP)
	return infraRsDomP, err
}

func (sm *ServiceManager) ReadDomain(tDn string, attachable_access_entity_profile string) (*models.Domain, error) {
	dn := fmt.Sprintf("uni/infra/attentp-%s/rsdomP-[%s]", attachable_access_entity_profile, tDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraRsDomP := models.DomainFromContainer(cont)
	return infraRsDomP, nil
}

func (sm *ServiceManager) DeleteDomain(tDn string, attachable_access_entity_profile string) error {
	dn := fmt.Sprintf("uni/infra/attentp-%s/rsdomP-[%s]", attachable_access_entity_profile, tDn)
	return sm.DeleteByDn(dn, models.InfrarsdompClassName)
}

func (sm *ServiceManager) UpdateDomain(tDn string, attachable_access_entity_profile string, description string, infraRsDomPattr models.DomainAttributes) (*models.Domain, error) {
	rn := fmt.Sprintf("rsdomP-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/infra/attentp-%s", attachable_access_entity_profile)
	infraRsDomP := models.NewDomain(rn, parentDn, description, infraRsDomPattr)

	infraRsDomP.Status = "modified"
	err := sm.Save(infraRsDomP)
	return infraRsDomP, err

}

func (sm *ServiceManager) ListDomain(attachable_access_entity_profile string) ([]*models.Domain, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/attentp-%s/infraRsDomP.json", baseurlStr, attachable_access_entity_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.DomainListFromContainer(cont)

	return list, err
}
