package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateOutofServiceFabricPath(tDn string, fabricRsOosPathAttr models.OutofServiceFabricPathAttributes) (*models.OutofServiceFabricPath, error) {
	rn := fmt.Sprintf(models.RnfabricRsOosPath, tDn)
	parentDn := fmt.Sprintf(models.ParentDnfabricRsOosPath)
	fabricRsOosPath := models.NewOutofServiceFabricPath(rn, parentDn, fabricRsOosPathAttr)
	err := sm.Save(fabricRsOosPath)
	return fabricRsOosPath, err
}

func (sm *ServiceManager) ReadOutofServiceFabricPath(tDn string) (*models.OutofServiceFabricPath, error) {
	dn := fmt.Sprintf(models.DnfabricRsOosPath, tDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	fabricRsOosPath := models.OutofServiceFabricPathFromContainer(cont)
	return fabricRsOosPath, nil
}

func (sm *ServiceManager) DeleteOutofServiceFabricPath(tDn string) error {
	dn := fmt.Sprintf(models.DnfabricRsOosPath, tDn)
	return sm.DeleteByDn(dn, models.FabricrsoospathClassName)
}

func (sm *ServiceManager) UpdateOutofServiceFabricPath(tDn string, fabricRsOosPathAttr models.OutofServiceFabricPathAttributes) (*models.OutofServiceFabricPath, error) {
	rn := fmt.Sprintf(models.RnfabricRsOosPath, tDn)
	parentDn := fmt.Sprintf(models.ParentDnfabricRsOosPath)
	fabricRsOosPath := models.NewOutofServiceFabricPath(rn, parentDn, fabricRsOosPathAttr)
	fabricRsOosPath.Status = "modified"
	err := sm.Save(fabricRsOosPath)
	return fabricRsOosPath, err
}

func (sm *ServiceManager) ListOutofServiceFabricPath() ([]*models.OutofServiceFabricPath, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/outofsvc/fabricRsOosPath.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.OutofServiceFabricPathListFromContainer(cont)
	return list, err
}
