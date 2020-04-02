package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateAccessSubPortBlock(name string, access_port_selector_type string, access_port_selector string, leaf_interface_profile string, description string, infraSubPortBlkattr models.AccessSubPortBlockAttributes) (*models.AccessSubPortBlock, error) {
	rn := fmt.Sprintf("subportblk-%s", name)
	parentDn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s", leaf_interface_profile, access_port_selector, access_port_selector_type)
	infraSubPortBlk := models.NewAccessSubPortBlock(rn, parentDn, description, infraSubPortBlkattr)
	err := sm.Save(infraSubPortBlk)
	return infraSubPortBlk, err
}

func (sm *ServiceManager) ReadAccessSubPortBlock(name string, access_port_selector_type string, access_port_selector string, leaf_interface_profile string) (*models.AccessSubPortBlock, error) {
	dn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s/subportblk-%s", leaf_interface_profile, access_port_selector, access_port_selector_type, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSubPortBlk := models.AccessSubPortBlockFromContainer(cont)
	return infraSubPortBlk, nil
}

func (sm *ServiceManager) DeleteAccessSubPortBlock(name string, access_port_selector_type string, access_port_selector string, leaf_interface_profile string) error {
	dn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s/subportblk-%s", leaf_interface_profile, access_port_selector, access_port_selector_type, name)
	return sm.DeleteByDn(dn, models.InfrasubportblkClassName)
}

func (sm *ServiceManager) UpdateAccessSubPortBlock(name string, access_port_selector_type string, access_port_selector string, leaf_interface_profile string, description string, infraSubPortBlkattr models.AccessSubPortBlockAttributes) (*models.AccessSubPortBlock, error) {
	rn := fmt.Sprintf("subportblk-%s", name)
	parentDn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s", leaf_interface_profile, access_port_selector, access_port_selector_type)
	infraSubPortBlk := models.NewAccessSubPortBlock(rn, parentDn, description, infraSubPortBlkattr)

	infraSubPortBlk.Status = "modified"
	err := sm.Save(infraSubPortBlk)
	return infraSubPortBlk, err

}

func (sm *ServiceManager) ListAccessSubPortBlock(access_port_selector_type string, access_port_selector string, leaf_interface_profile string) ([]*models.AccessSubPortBlock, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/accportprof-%s/hports-%s-typ-%s/infraSubPortBlk.json", baseurlStr, leaf_interface_profile, access_port_selector, access_port_selector_type)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.AccessSubPortBlockListFromContainer(cont)

	return list, err
}
