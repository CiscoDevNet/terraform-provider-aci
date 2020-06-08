package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCtrlrMaintP(description string, maintCtrlrMaintPAttr models.CtrlrMaintPAttributes) (*models.CtrlrMaintP, error) {
	rn := fmt.Sprintf("ctrlrmaintpol")
	parentDn := fmt.Sprintf("uni/controller")
	maintCtrlrMaintP := models.NewCtrlrMaintP(rn, parentDn, description, maintCtrlrMaintPAttr)
	err := sm.Save(maintCtrlrMaintP)
	return maintCtrlrMaintP, err
}

func (sm *ServiceManager) ReadCtrlrMaintP(name string) (*models.CtrlrMaintP, error) {
	dn := fmt.Sprintf("uni/controller/ctrlrmaintpol")
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	maintCtrlrMaintP := models.CtrlrMaintPFromContainer(cont)
	return maintCtrlrMaintP, nil
}

func (sm *ServiceManager) ListCtrlrMaintP() ([]*models.CtrlrMaintP, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/maintCtrlrMaintP.json", baseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.CtrlrMaintPListFromContainer(cont)
	return list, err
}
