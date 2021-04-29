package client

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) createFabricBlkParentDn(parentMOType string, parentMOName string) (string, error) {
	var parentMOPrefix string
	switch parentMOType {
	case "maintMaintGrp":
		parentMOPrefix = "maintgrp"
	case "firmwareFwGrp":
		parentMOPrefix = "fwgrp"
	default:
		return "", fmt.Errorf("Unsupported parentMO Type: %s", parentMOType)
	}

	if parentMOName == "" {
		return "", fmt.Errorf("parentMO Name missing")
	}

	// uni/fabric/fwgrp-myFwGrp
	// uni/fabric/maintgrp-myMaintGrp
	parentDn := fmt.Sprintf("uni/fabric/%s-%s", parentMOPrefix, parentMOName)
	return parentDn, nil
}

func (sm *ServiceManager) CreateNodeBlk(parentMOType string, parentMOName string, name string, description string, fabricNodeBlkAttr models.NodeBlkAttributes) (*models.NodeBlk, error) {
	rn := fmt.Sprintf("nodeblk-%s", name)
	parentDn, err := sm.createFabricBlkParentDn(parentMOType, parentMOName)
	if err != nil {
		return nil, err
	}

	fabricNodeBlk := models.NewNodeBlk(rn, parentDn, description, fabricNodeBlkAttr)
	err = sm.Save(fabricNodeBlk)
	return fabricNodeBlk, err
}

func (sm *ServiceManager) ReadNodeBlk(parentMOType string, parentMOName string, name string) (*models.NodeBlk, error) {
	parentDn, err := sm.createFabricBlkParentDn(parentMOType, parentMOName)
	if err != nil {
		return nil, err
	}
	dn := fmt.Sprintf("%s/nodeblk-%s", parentDn, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	fabricNodeBlk := models.NodeBlkFromContainer(cont)
	return fabricNodeBlk, nil
}

func (sm *ServiceManager) DeleteNodeBlk(parentMOType string, parentMOName string, name string) error {
	parentDn, err := sm.createFabricBlkParentDn(parentMOType, parentMOName)
	if err != nil {
		return err
	}
	dn := fmt.Sprintf("%s/nodeblk-%s", parentDn, name)
	return sm.DeleteByDn(dn, models.FabricNodeBlkClassName)
}

func (sm *ServiceManager) UpdateNodeBlk(parentMOType string, parentMOName string, name string, description string, fabricNodeBlkAttr models.NodeBlkAttributes) (*models.NodeBlk, error) {
	rn := fmt.Sprintf("nodeblk-%s", name)
	parentDn, err := sm.createFabricBlkParentDn(parentMOType, parentMOName)
	if err != nil {
		return nil, err
	}
	fabricNodeBlk := models.NewNodeBlk(rn, parentDn, description, fabricNodeBlkAttr)
	fabricNodeBlk.Status = "modified"
	err = sm.Save(fabricNodeBlk)
	return fabricNodeBlk, err

}

func (sm *ServiceManager) ListNodeBlk() ([]*models.NodeBlk, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/fabricNodeBlk.json", baseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.NodeBlkListFromContainer(cont)
	return list, err
}
