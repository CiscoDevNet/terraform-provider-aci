package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateNodeBlockFW(name string, firmware_group string, description string, fabricNodeBlkattr models.NodeBlockAttributesFW) (*models.NodeBlockFW, error) {
	rn := fmt.Sprintf("nodeblk-%s", name)
	parentDn := fmt.Sprintf("uni/fabric/fwgrp-%s", firmware_group)
	fabricNodeBlk := models.NewNodeBlockFW(rn, parentDn, description, fabricNodeBlkattr)
	err := sm.Save(fabricNodeBlk)
	return fabricNodeBlk, err
}

func (sm *ServiceManager) ReadNodeBlockFW(name string, firmware_group string) (*models.NodeBlockFW, error) {
	dn := fmt.Sprintf("uni/fabric/fwgrp-%s/nodeblk-%s", firmware_group, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricNodeBlk := models.NodeBlockFromContainer(cont)
	return fabricNodeBlk, nil
}

func (sm *ServiceManager) DeleteNodeBlockFW(name string, firmware_group string) error {
	dn := fmt.Sprintf("uni/fabric/fwgrp-%s/nodeblk-%s", firmware_group, name)
	return sm.DeleteByDn(dn, models.FabricnodeblkClassNameFW)
}

func (sm *ServiceManager) UpdateNodeBlockFW(name string, firmware_group string, description string, fabricNodeBlkattr models.NodeBlockAttributesFW) (*models.NodeBlockFW, error) {
	rn := fmt.Sprintf("nodeblk-%s", name)
	parentDn := fmt.Sprintf("uni/fabric/fwgrp-%s", firmware_group)
	fabricNodeBlk := models.NewNodeBlockFW(rn, parentDn, description, fabricNodeBlkattr)

	fabricNodeBlk.Status = "modified"
	err := sm.Save(fabricNodeBlk)
	return fabricNodeBlk, err

}

func (sm *ServiceManager) ListNodeBlockFW(firmware_group string) ([]*models.NodeBlockFW, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fabric/fwgrp-%s/fabricNodeBlk.json", baseurlStr, firmware_group)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.NodeBlockListFromContainer(cont)

	return list, err
}
