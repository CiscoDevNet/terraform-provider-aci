package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSpineInterfaceProfile(name string, description string, infraSpAccPortPattr models.SpineInterfaceProfileAttributes) (*models.SpineInterfaceProfile, error) {
	rn := fmt.Sprintf("infra/spaccportprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraSpAccPortP := models.NewSpineInterfaceProfile(rn, parentDn, description, infraSpAccPortPattr)
	err := sm.Save(infraSpAccPortP)
	return infraSpAccPortP, err
}

func (sm *ServiceManager) ReadSpineInterfaceProfile(name string) (*models.SpineInterfaceProfile, error) {
	dn := fmt.Sprintf("uni/infra/spaccportprof-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpAccPortP := models.SpineInterfaceProfileFromContainer(cont)
	return infraSpAccPortP, nil
}

func (sm *ServiceManager) DeleteSpineInterfaceProfile(name string) error {
	dn := fmt.Sprintf("uni/infra/spaccportprof-%s", name)
	return sm.DeleteByDn(dn, models.InfraspaccportpClassName)
}

func (sm *ServiceManager) UpdateSpineInterfaceProfile(name string, description string, infraSpAccPortPattr models.SpineInterfaceProfileAttributes) (*models.SpineInterfaceProfile, error) {
	rn := fmt.Sprintf("infra/spaccportprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraSpAccPortP := models.NewSpineInterfaceProfile(rn, parentDn, description, infraSpAccPortPattr)

	infraSpAccPortP.Status = "modified"
	err := sm.Save(infraSpAccPortP)
	return infraSpAccPortP, err

}

func (sm *ServiceManager) ListSpineInterfaceProfile() ([]*models.SpineInterfaceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infraSpAccPortP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.SpineInterfaceProfileListFromContainer(cont)

	return list, err
}
