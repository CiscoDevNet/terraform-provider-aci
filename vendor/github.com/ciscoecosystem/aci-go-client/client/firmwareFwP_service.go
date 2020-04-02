package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFirmwarePolicy(name string, description string, firmwareFwPattr models.FirmwarePolicyAttributes) (*models.FirmwarePolicy, error) {
	rn := fmt.Sprintf("fabric/fwpol-%s", name)
	parentDn := fmt.Sprintf("uni")
	firmwareFwP := models.NewFirmwarePolicy(rn, parentDn, description, firmwareFwPattr)
	err := sm.Save(firmwareFwP)
	return firmwareFwP, err
}

func (sm *ServiceManager) ReadFirmwarePolicy(name string) (*models.FirmwarePolicy, error) {
	dn := fmt.Sprintf("uni/fabric/fwpol-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	firmwareFwP := models.FirmwarePolicyFromContainer(cont)
	return firmwareFwP, nil
}

func (sm *ServiceManager) DeleteFirmwarePolicy(name string) error {
	dn := fmt.Sprintf("uni/fabric/fwpol-%s", name)
	return sm.DeleteByDn(dn, models.FirmwarefwpClassName)
}

func (sm *ServiceManager) UpdateFirmwarePolicy(name string, description string, firmwareFwPattr models.FirmwarePolicyAttributes) (*models.FirmwarePolicy, error) {
	rn := fmt.Sprintf("fabric/fwpol-%s", name)
	parentDn := fmt.Sprintf("uni")
	firmwareFwP := models.NewFirmwarePolicy(rn, parentDn, description, firmwareFwPattr)

	firmwareFwP.Status = "modified"
	err := sm.Save(firmwareFwP)
	return firmwareFwP, err

}

func (sm *ServiceManager) ListFirmwarePolicy() ([]*models.FirmwarePolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/firmwareFwP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FirmwarePolicyListFromContainer(cont)

	return list, err
}
