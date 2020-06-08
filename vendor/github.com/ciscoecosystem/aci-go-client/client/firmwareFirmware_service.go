package client

import (
        "fmt"
        "github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFirmware(name string, description string, firmwareFirmwareAttr models.FirmwareAttributes) (*models.Firmware, error) {
        rn := fmt.Sprintf("fw-%s",name)
        parentDn := fmt.Sprintf("uni/fabric/fwrepop")
        firmwareFirmware := models.NewFirmware(rn, parentDn, description, firmwareFirmwareAttr)
        err := sm.Save(firmwareFirmware)
        return firmwareFirmware, err
}

func (sm *ServiceManager) ReadFirmware(name string) (*models.Firmware, error) {
        dn := fmt.Sprintf("fwrepop/fw-%s", name )
        cont, err := sm.Get(dn)
        if err != nil {
                return nil, err
        }
        firmwareFirmware := models.FirmwareFromContainer(cont)
        return firmwareFirmware, nil
}

func (sm *ServiceManager) DeleteFirmware(name string ) error {
        dn := fmt.Sprintf("uni/fabric/fwrepop/osrc-%s", name)
        return sm.DeleteByDn(dn, models.FirmwareFirmwareClassName)
}

func (sm *ServiceManager) UpdateFirmware(name string, description string, firmwareFirmwareAttr models.FirmwareAttributes) (*models.Firmware, error) {
        rn := fmt.Sprintf("osrc-%s",name)
        parentDn := fmt.Sprintf("uni/fabric/fwrepop")
        firmwareFirmware := models.NewFirmware(rn, parentDn, description, firmwareFirmwareAttr)
        firmwareFirmware.Status = "modified"
        err := sm.Save(firmwareFirmware)
        return firmwareFirmware, err

}

func (sm *ServiceManager) ListFirmware() ([]*models.Firmware, error) {
        baseurlStr := "/api/node/class"
        dnUrl := fmt.Sprintf("%s/firmwareFirmware.json", baseurlStr )
        cont, err := sm.GetViaURL(dnUrl)
        list := models.FirmwareListFromContainer(cont)
        return list, err
}

