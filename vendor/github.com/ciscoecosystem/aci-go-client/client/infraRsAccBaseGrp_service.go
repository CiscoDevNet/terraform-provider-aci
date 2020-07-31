package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateAccessAccessGroup(access_port_selector_type string, access_port_selector string, leaf_interface_profile string, description string, infraRsAccBaseGrpattr models.AccessAccessGroupAttributes) (*models.AccessAccessGroup, error) {
	rn := fmt.Sprintf("rsaccBaseGrp")
	parentDn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s", leaf_interface_profile, access_port_selector, access_port_selector_type)
	infraRsAccBaseGrp := models.NewAccessAccessGroup(rn, parentDn, description, infraRsAccBaseGrpattr)
	err := sm.Save(infraRsAccBaseGrp)
	return infraRsAccBaseGrp, err
}

func (sm *ServiceManager) ReadAccessAccessGroup(access_port_selector_type string, access_port_selector string, leaf_interface_profile string) (*models.AccessAccessGroup, error) {
	dn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s/rsaccBaseGrp", leaf_interface_profile, access_port_selector, access_port_selector_type)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraRsAccBaseGrp := models.AccessAccessGroupFromContainer(cont)
	return infraRsAccBaseGrp, nil
}

func (sm *ServiceManager) DeleteAccessAccessGroup(access_port_selector_type string, access_port_selector string, leaf_interface_profile string) error {
	dn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s/rsaccBaseGrp", leaf_interface_profile, access_port_selector, access_port_selector_type)
	return sm.DeleteByDn(dn, models.InfrarsaccbasegrpClassName)
}

func (sm *ServiceManager) UpdateAccessAccessGroup(access_port_selector_type string, access_port_selector string, leaf_interface_profile string, description string, infraRsAccBaseGrpattr models.AccessAccessGroupAttributes) (*models.AccessAccessGroup, error) {
	rn := fmt.Sprintf("rsaccBaseGrp")
	parentDn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s", leaf_interface_profile, access_port_selector, access_port_selector_type)
	infraRsAccBaseGrp := models.NewAccessAccessGroup(rn, parentDn, description, infraRsAccBaseGrpattr)

	infraRsAccBaseGrp.Status = "modified"
	err := sm.Save(infraRsAccBaseGrp)
	return infraRsAccBaseGrp, err

}

func (sm *ServiceManager) ListAccessAccessGroup(access_port_selector_type string, access_port_selector string, leaf_interface_profile string) ([]*models.AccessAccessGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/accportprof-%s/hports-%s-typ-%s/infraRsAccBaseGrp.json", baseurlStr, leaf_interface_profile, access_port_selector, access_port_selector_type)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.AccessAccessGroupListFromContainer(cont)

	return list, err
}
