package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLeafInterfaceProfile(name string, description string, infraAccPortPattr models.LeafInterfaceProfileAttributes) (*models.LeafInterfaceProfile, error) {
	rn := fmt.Sprintf("infra/accportprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraAccPortP := models.NewLeafInterfaceProfile(rn, parentDn, description, infraAccPortPattr)
	err := sm.Save(infraAccPortP)
	return infraAccPortP, err
}

func (sm *ServiceManager) ReadLeafInterfaceProfile(name string) (*models.LeafInterfaceProfile, error) {
	dn := fmt.Sprintf("uni/infra/accportprof-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraAccPortP := models.LeafInterfaceProfileFromContainer(cont)
	return infraAccPortP, nil
}

func (sm *ServiceManager) DeleteLeafInterfaceProfile(name string) error {
	dn := fmt.Sprintf("uni/infra/accportprof-%s", name)
	return sm.DeleteByDn(dn, models.InfraaccportpClassName)
}

func (sm *ServiceManager) UpdateLeafInterfaceProfile(name string, description string, infraAccPortPattr models.LeafInterfaceProfileAttributes) (*models.LeafInterfaceProfile, error) {
	rn := fmt.Sprintf("infra/accportprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraAccPortP := models.NewLeafInterfaceProfile(rn, parentDn, description, infraAccPortPattr)

	infraAccPortP.Status = "modified"
	err := sm.Save(infraAccPortP)
	return infraAccPortP, err

}

func (sm *ServiceManager) ListLeafInterfaceProfile() ([]*models.LeafInterfaceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infraAccPortP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LeafInterfaceProfileListFromContainer(cont)

	return list, err
}
