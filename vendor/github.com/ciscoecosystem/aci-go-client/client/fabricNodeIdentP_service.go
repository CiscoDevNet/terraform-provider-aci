package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFabricNodeMember(serial string, description string, fabricNodeIdentPattr models.FabricNodeMemberAttributes) (*models.FabricNodeMember, error) {
	rn := fmt.Sprintf("controller/nodeidentpol/nodep-%s", serial)
	parentDn := fmt.Sprintf("uni")
	fabricNodeIdentP := models.NewFabricNodeMember(rn, parentDn, description, fabricNodeIdentPattr)
	err := sm.Save(fabricNodeIdentP)
	return fabricNodeIdentP, err
}

func (sm *ServiceManager) ReadFabricNodeMember(serial string) (*models.FabricNodeMember, error) {
	dn := fmt.Sprintf("uni/controller/nodeidentpol/nodep-%s", serial)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricNodeIdentP := models.FabricNodeMemberFromContainer(cont)
	return fabricNodeIdentP, nil
}

func (sm *ServiceManager) DeleteFabricNodeMember(serial string) error {
	dn := fmt.Sprintf("uni/controller/nodeidentpol/nodep-%s", serial)
	return sm.DeleteByDn(dn, models.FabricnodeidentpClassName)
}

func (sm *ServiceManager) UpdateFabricNodeMember(serial string, description string, fabricNodeIdentPattr models.FabricNodeMemberAttributes) (*models.FabricNodeMember, error) {
	rn := fmt.Sprintf("controller/nodeidentpol/nodep-%s", serial)
	parentDn := fmt.Sprintf("uni")
	fabricNodeIdentP := models.NewFabricNodeMember(rn, parentDn, description, fabricNodeIdentPattr)

	fabricNodeIdentP.Status = "modified"
	err := sm.Save(fabricNodeIdentP)
	return fabricNodeIdentP, err

}

func (sm *ServiceManager) ListFabricNodeMember() ([]*models.FabricNodeMember, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fabricNodeIdentP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FabricNodeMemberListFromContainer(cont)

	return list, err
}
