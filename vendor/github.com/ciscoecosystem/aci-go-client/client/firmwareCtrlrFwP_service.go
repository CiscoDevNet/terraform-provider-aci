package client

import (
        "fmt"
        "github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCtrlrFwP(description string, firmwareCtrlrFwPAttr models.CtrlrFwPAttributes) (*models.CtrlrFwP, error) {
        rn := "ctrlrfwpol"
        parentDn := fmt.Sprintf("uni/controller")
        firmwareCtrlrFwP := models.NewCtrlrFwP(rn, parentDn, description, firmwareCtrlrFwPAttr)
        err := sm.Save(firmwareCtrlrFwP)
        return firmwareCtrlrFwP, err
}

func (sm *ServiceManager) ReadCtrlrFwP(name string) (*models.CtrlrFwP, error) {
        dn := fmt.Sprintf("uni/controller/ctrlrfwpol")
        cont, err := sm.Get(dn)
        if err != nil {
                return nil, err
        }
        firmwareCtrlrFwP := models.CtrlrFwPFromContainer(cont)
        return firmwareCtrlrFwP, nil
}

func (sm *ServiceManager) DeleteCtrlrFwP() error {
        dn := fmt.Sprintf("uni/controller/ctrlrfwpol")
        return sm.DeleteByDn(dn, models.FirmwareCtrlrFwPClassName)
}

func (sm *ServiceManager) UpdateCtrlrFwP(description string, firmwareCtrlrFwPAttr models.CtrlrFwPAttributes) (*models.CtrlrFwP, error) {
        rn := "ctrlrfwpol"
        parentDn := fmt.Sprintf("uni/controller")
        firmwareCtrlrFwP := models.NewCtrlrFwP(rn, parentDn, description, firmwareCtrlrFwPAttr)
        firmwareCtrlrFwP.Status = "modified"
        err := sm.Save(firmwareCtrlrFwP)
        return firmwareCtrlrFwP, err

}

func (sm *ServiceManager) ListCtrlrFwP() ([]*models.CtrlrFwP, error) {
        baseurlStr := "/api/node/class"
        dnUrl := fmt.Sprintf("%s/firmwareCtrlrFwP.json", baseurlStr )
        cont, err := sm.GetViaURL(dnUrl)
        list := models.CtrlrFwPListFromContainer(cont)
        return list, err
}

