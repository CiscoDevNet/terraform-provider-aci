package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFirmwareGroup(name string, description string, firmwareFwGrpattr models.FirmwareGroupAttributes) (*models.FirmwareGroup, error) {
	rn := fmt.Sprintf("fabric/fwgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	firmwareFwGrp := models.NewFirmwareGroup(rn, parentDn, description, firmwareFwGrpattr)
	err := sm.Save(firmwareFwGrp)
	return firmwareFwGrp, err
}

func (sm *ServiceManager) ReadFirmwareGroup(name string) (*models.FirmwareGroup, error) {
	dn := fmt.Sprintf("uni/fabric/fwgrp-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	firmwareFwGrp := models.FirmwareGroupFromContainer(cont)
	return firmwareFwGrp, nil
}

func (sm *ServiceManager) DeleteFirmwareGroup(name string) error {
	dn := fmt.Sprintf("uni/fabric/fwgrp-%s", name)
	return sm.DeleteByDn(dn, models.FirmwarefwgrpClassName)
}

func (sm *ServiceManager) UpdateFirmwareGroup(name string, description string, firmwareFwGrpattr models.FirmwareGroupAttributes) (*models.FirmwareGroup, error) {
	rn := fmt.Sprintf("fabric/fwgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	firmwareFwGrp := models.NewFirmwareGroup(rn, parentDn, description, firmwareFwGrpattr)

	firmwareFwGrp.Status = "modified"
	err := sm.Save(firmwareFwGrp)
	return firmwareFwGrp, err

}

func (sm *ServiceManager) ListFirmwareGroup() ([]*models.FirmwareGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/firmwareFwGrp.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FirmwareGroupListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationfirmwareRsFwgrppFromFirmwareGroup(parentDn, tnFirmwareFwPName string) error {
	dn := fmt.Sprintf("%s/rsfwgrpp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFirmwareFwPName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "firmwareRsFwgrpp", dn, tnFirmwareFwPName))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) ReadRelationfirmwareRsFwgrppFromFirmwareGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "firmwareRsFwgrpp")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "firmwareRsFwgrpp")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
