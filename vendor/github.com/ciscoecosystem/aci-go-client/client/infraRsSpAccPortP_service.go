package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateInterfaceProfile(tDn string, spine_profile string, description string, infraRsSpAccPortPattr models.InterfaceProfileAttributes) (*models.InterfaceProfile, error) {
	rn := fmt.Sprintf("rsspAccPortP-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/infra/spprof-%s", spine_profile)
	infraRsSpAccPortP := models.NewInterfaceProfile(rn, parentDn, description, infraRsSpAccPortPattr)
	err := sm.Save(infraRsSpAccPortP)
	return infraRsSpAccPortP, err
}

func (sm *ServiceManager) ReadInterfaceProfile(tDn string, spine_profile string) (*models.InterfaceProfile, error) {
	dn := fmt.Sprintf("uni/infra/spprof-%s/rsspAccPortP-[%s]", spine_profile, tDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraRsSpAccPortP := models.InterfaceProfileFromContainer(cont)
	return infraRsSpAccPortP, nil
}

func (sm *ServiceManager) DeleteInterfaceProfile(tDn string, spine_profile string) error {
	dn := fmt.Sprintf("uni/infra/spprof-%s/rsspAccPortP-[%s]", spine_profile, tDn)
	return sm.DeleteByDn(dn, models.InfrarsspaccportpClassName)
}

func (sm *ServiceManager) UpdateInterfaceProfile(tDn string, spine_profile string, description string, infraRsSpAccPortPattr models.InterfaceProfileAttributes) (*models.InterfaceProfile, error) {
	rn := fmt.Sprintf("rsspAccPortP-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/infra/spprof-%s", spine_profile)
	infraRsSpAccPortP := models.NewInterfaceProfile(rn, parentDn, description, infraRsSpAccPortPattr)

	infraRsSpAccPortP.Status = "modified"
	err := sm.Save(infraRsSpAccPortP)
	return infraRsSpAccPortP, err

}

func (sm *ServiceManager) ListInterfaceProfile(spine_profile string) ([]*models.InterfaceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/spprof-%s/infraRsSpAccPortP.json", baseurlStr, spine_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.InterfaceProfileListFromContainer(cont)

	return list, err
}
