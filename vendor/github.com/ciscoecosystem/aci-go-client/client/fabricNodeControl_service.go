package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFabricNodeControl(name string, description string, nameAlias string, fabricNodeControlAttr models.FabricNodeControlAttributes) (*models.FabricNodeControl, error) {
	rn := fmt.Sprintf(models.RnfabricNodeControl, name)
	parentDn := fmt.Sprintf(models.ParentDnfabricNodeControl)
	fabricNodeControl := models.NewFabricNodeControl(rn, parentDn, description, nameAlias, fabricNodeControlAttr)
	err := sm.Save(fabricNodeControl)
	return fabricNodeControl, err
}

func (sm *ServiceManager) ReadFabricNodeControl(name string) (*models.FabricNodeControl, error) {
	dn := fmt.Sprintf(models.DnfabricNodeControl, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	fabricNodeControl := models.FabricNodeControlFromContainer(cont)
	return fabricNodeControl, nil
}

func (sm *ServiceManager) DeleteFabricNodeControl(name string) error {
	dn := fmt.Sprintf(models.DnfabricNodeControl, name)
	return sm.DeleteByDn(dn, models.FabricnodecontrolClassName)
}

func (sm *ServiceManager) UpdateFabricNodeControl(name string, description string, nameAlias string, fabricNodeControlAttr models.FabricNodeControlAttributes) (*models.FabricNodeControl, error) {
	rn := fmt.Sprintf(models.RnfabricNodeControl, name)
	parentDn := fmt.Sprintf(models.ParentDnfabricNodeControl)
	fabricNodeControl := models.NewFabricNodeControl(rn, parentDn, description, nameAlias, fabricNodeControlAttr)
	fabricNodeControl.Status = "modified"
	err := sm.Save(fabricNodeControl)
	return fabricNodeControl, err
}

func (sm *ServiceManager) ListFabricNodeControl() ([]*models.FabricNodeControl, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/fabricNodeControl.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.FabricNodeControlListFromContainer(cont)
	return list, err
}
