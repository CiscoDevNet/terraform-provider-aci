package client

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) ReadFirmwareCtrlrRunning(podId uint32, nodeId uint32) (*models.FirmwareCtrlrRunning, error) {
	dn := fmt.Sprintf("topology/pod-%d/node-%d/sys/ctrlrfwstatuscont/ctrlrrunning", podId, nodeId)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	firmwareFirmwareCtrlrRunning := models.FirmwareCtrlrRunningFromContainer(cont)
	return firmwareFirmwareCtrlrRunning, nil
}

func (sm *ServiceManager) ListFirmwareCtrlrRunning() ([]*models.FirmwareCtrlrRunning, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/firmwareCtrlrRunning.json", baseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.FirmwareCtrlrRunningListFromContainer(cont)
	return list, err
}
