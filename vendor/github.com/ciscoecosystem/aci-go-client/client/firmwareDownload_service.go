package client

import (
        "fmt"
        "github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateDownload(name string, description string, firmwareDownloadAttr models.DownloadAttributes) (*models.Download, error) {
        rn := fmt.Sprintf("download-%s",name)
        parentDn := fmt.Sprintf("uni/fabric/fwrepop")
        firmwareDownload := models.NewDownload(rn, parentDn, description, firmwareDownloadAttr)
        err := sm.Save(firmwareDownload)
        return firmwareDownload, err
}

func (sm *ServiceManager) ReadDownload(name string) (*models.Download, error) {
        dn := fmt.Sprintf("uni/fabric/fwrepop/download-%s", name )
        cont, err := sm.Get(dn)
        if err != nil {
                return nil, err
        }
        firmwareDownload := models.DownloadFromContainer(cont)
        return firmwareDownload, nil
}

func (sm *ServiceManager) DeleteDownload(name string ) error {
        dn := fmt.Sprintf("uni/fabric/fwrepop/download-%s", name)
        return sm.DeleteByDn(dn, models.FirmwareDownloadClassName)
}

func (sm *ServiceManager) UpdateDownload(name string, description string, firmwareDownloadAttr models.DownloadAttributes) (*models.Download, error) {
        rn := fmt.Sprintf("download-%s",name)
        parentDn := fmt.Sprintf("uni/fabric/fwrepop")
        firmwareDownload := models.NewDownload(rn, parentDn, description, firmwareDownloadAttr)
        firmwareDownload.Status = "modified"
        err := sm.Save(firmwareDownload)
        return firmwareDownload, err

}

func (sm *ServiceManager) ListDownload() ([]*models.Download, error) {
        baseurlStr := "/api/node/class"
        dnUrl := fmt.Sprintf("%s/firmwareDownload.json", baseurlStr )
        cont, err := sm.GetViaURL(dnUrl)
        list := models.DownloadListFromContainer(cont)
        return list, err
}

