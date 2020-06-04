package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFirmwareDownloadTask(name string, description string, firmwareOSourceattr models.FirmwareDownloadTaskAttributes) (*models.FirmwareDownloadTask, error) {
	rn := fmt.Sprintf("fabric/fwrepop/osrc-%s", name)
	parentDn := fmt.Sprintf("uni")
	firmwareOSource := models.NewFirmwareDownloadTask(rn, parentDn, description, firmwareOSourceattr)
	err := sm.Save(firmwareOSource)
	return firmwareOSource, err
}

func (sm *ServiceManager) ReadFirmwareDownloadTask(name string) (*models.FirmwareDownloadTask, error) {
	dn := fmt.Sprintf("uni/fabric/fwrepop/osrc-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	firmwareOSource := models.FirmwareDownloadTaskFromContainer(cont)
	return firmwareOSource, nil
}

func (sm *ServiceManager) DeleteFirmwareDownloadTask(name string) error {
	dn := fmt.Sprintf("uni/fabric/fwrepop/osrc-%s", name)
	return sm.DeleteByDn(dn, models.FirmwareosourceClassName)
}

func (sm *ServiceManager) UpdateFirmwareDownloadTask(name string, description string, firmwareOSourceattr models.FirmwareDownloadTaskAttributes) (*models.FirmwareDownloadTask, error) {
	rn := fmt.Sprintf("fabric/fwrepop/osrc-%s", name)
	parentDn := fmt.Sprintf("uni")
	firmwareOSource := models.NewFirmwareDownloadTask(rn, parentDn, description, firmwareOSourceattr)

	firmwareOSource.Status = "modified"
	err := sm.Save(firmwareOSource)
	return firmwareOSource, err

}

func (sm *ServiceManager) ListFirmwareDownloadTask() ([]*models.FirmwareDownloadTask, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/firmwareOSource.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FirmwareDownloadTaskListFromContainer(cont)

	return list, err
}

// START: Variable/Struct/Fuction Naming per ACI SDK Model Definitions
func (sm *ServiceManager) CreateOSource(name string, description string, firmwareOSourceAttr models.OSourceAttributes) (*models.OSource, error) {
	rn := fmt.Sprintf("osrc-%s", name)
	parentDn := fmt.Sprintf("uni/fabric/fwrepop")
	firmwareOSource := models.NewOSource(rn, parentDn, description, firmwareOSourceAttr)
	err := sm.Save(firmwareOSource)
	return firmwareOSource, err
}

func (sm *ServiceManager) ReadOSource(name string) (*models.OSource, error) {
	dn := fmt.Sprintf("uni/fabric/fwrepop/osrc-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	firmwareOSource := models.OSourceFromContainer(cont)
	return firmwareOSource, nil
}

func (sm *ServiceManager) DeleteOSource(name string) error {
	dn := fmt.Sprintf("uni/fabric/fwrepop/osrc-%s", name)
	return sm.DeleteByDn(dn, models.FirmwareOSourceClassName)
}

func (sm *ServiceManager) UpdateOSource(name string, description string, firmwareOSourceAttr models.OSourceAttributes) (*models.OSource, error) {
	rn := fmt.Sprintf("osrc-%s", name)
	parentDn := fmt.Sprintf("uni/fabric/fwrepop")
	firmwareOSource := models.NewOSource(rn, parentDn, description, firmwareOSourceAttr)
	firmwareOSource.Status = "modified"
	err := sm.Save(firmwareOSource)
	return firmwareOSource, err

}

func (sm *ServiceManager) ListOSource() ([]*models.OSource, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/firmwareOSource.json", baseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.OSourceListFromContainer(cont)
	return list, err
}
