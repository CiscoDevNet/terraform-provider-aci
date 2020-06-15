package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateNodeBlock(name string, switch_association_type string, switch_association string, leaf_profile string, description string, infraNodeBlkattr models.NodeBlockAttributes) (*models.NodeBlock, error) {
	rn := fmt.Sprintf("nodeblk-%s", name)
	parentDn := fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-%s", leaf_profile, switch_association, switch_association_type)
	infraNodeBlk := models.NewNodeBlock(rn, parentDn, description, infraNodeBlkattr)
	err := sm.Save(infraNodeBlk)
	return infraNodeBlk, err
}

func (sm *ServiceManager) ReadNodeBlock(name string, switch_association_type string, switch_association string, leaf_profile string) (*models.NodeBlock, error) {
	dn := fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-%s/nodeblk-%s", leaf_profile, switch_association, switch_association_type, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraNodeBlk := models.NodeBlockFromContainerBLK(cont)
	return infraNodeBlk, nil
}

func (sm *ServiceManager) DeleteNodeBlock(name string, switch_association_type string, switch_association string, leaf_profile string) error {
	dn := fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-%s/nodeblk-%s", leaf_profile, switch_association, switch_association_type, name)
	return sm.DeleteByDn(dn, models.InfranodeblkClassName)
}

func (sm *ServiceManager) UpdateNodeBlock(name string, switch_association_type string, switch_association string, leaf_profile string, description string, infraNodeBlkattr models.NodeBlockAttributes) (*models.NodeBlock, error) {
	rn := fmt.Sprintf("nodeblk-%s", name)
	parentDn := fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-%s", leaf_profile, switch_association, switch_association_type)
	infraNodeBlk := models.NewNodeBlock(rn, parentDn, description, infraNodeBlkattr)

	infraNodeBlk.Status = "modified"
	err := sm.Save(infraNodeBlk)
	return infraNodeBlk, err

}

func (sm *ServiceManager) ListNodeBlock(switch_association_type string, switch_association string, leaf_profile string) ([]*models.NodeBlock, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/nprof-%s/leaves-%s-typ-%s/infraNodeBlk.json", baseurlStr, leaf_profile, switch_association, switch_association_type)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.NodeBlockListFromContainerBLK(cont)

	return list, err
}
