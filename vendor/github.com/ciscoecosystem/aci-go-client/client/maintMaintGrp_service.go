package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreatePODMaintenanceGroup(name string, description string, maintMaintGrpattr models.PODMaintenanceGroupAttributes) (*models.PODMaintenanceGroup, error) {
	rn := fmt.Sprintf("fabric/maintgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	maintMaintGrp := models.NewPODMaintenanceGroup(rn, parentDn, description, maintMaintGrpattr)
	err := sm.Save(maintMaintGrp)
	return maintMaintGrp, err
}

func (sm *ServiceManager) ReadPODMaintenanceGroup(name string) (*models.PODMaintenanceGroup, error) {
	dn := fmt.Sprintf("uni/fabric/maintgrp-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	maintMaintGrp := models.PODMaintenanceGroupFromContainer(cont)
	return maintMaintGrp, nil
}

func (sm *ServiceManager) DeletePODMaintenanceGroup(name string) error {
	dn := fmt.Sprintf("uni/fabric/maintgrp-%s", name)
	return sm.DeleteByDn(dn, models.MaintmaintgrpClassName)
}

func (sm *ServiceManager) UpdatePODMaintenanceGroup(name string, description string, maintMaintGrpattr models.PODMaintenanceGroupAttributes) (*models.PODMaintenanceGroup, error) {
	rn := fmt.Sprintf("fabric/maintgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	maintMaintGrp := models.NewPODMaintenanceGroup(rn, parentDn, description, maintMaintGrpattr)

	maintMaintGrp.Status = "modified"
	err := sm.Save(maintMaintGrp)
	return maintMaintGrp, err

}

func (sm *ServiceManager) ListPODMaintenanceGroup() ([]*models.PODMaintenanceGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/maintMaintGrp.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.PODMaintenanceGroupListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationmaintRsMgrppFromPODMaintenanceGroup(parentDn, tnMaintMaintPName string) error {
	dn := fmt.Sprintf("%s/rsmgrpp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMaintMaintPName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "maintRsMgrpp", dn, tnMaintMaintPName))

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

func (sm *ServiceManager) ReadRelationmaintRsMgrppFromPODMaintenanceGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "maintRsMgrpp")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "maintRsMgrpp")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

// START: Variable/Struct/Fuction Naming per ACI SDK Model Definitions
func (sm *ServiceManager) CreateMaintGrp(name string, description string, maintMaintGrpAttr models.MaintGrpAttributes) (*models.MaintGrp, error) {
	rn := fmt.Sprintf("maintgrp-%s", name)
	parentDn := fmt.Sprintf("uni/fabric")
	maintMaintGrp := models.NewMaintGrp(rn, parentDn, description, maintMaintGrpAttr)
	err := sm.Save(maintMaintGrp)
	return maintMaintGrp, err
}

func (sm *ServiceManager) ReadMaintGrp(name string) (*models.MaintGrp, error) {
	dn := fmt.Sprintf("uni/fabric/maintgrp-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	maintMaintGrp := models.MaintGrpFromContainer(cont)
	return maintMaintGrp, nil
}

func (sm *ServiceManager) DeleteMaintGrp(name string) error {
	dn := fmt.Sprintf("uni/fabric/maintgrp-%s", name)
	return sm.DeleteByDn(dn, models.MaintMaintGrpClassName)
}

func (sm *ServiceManager) UpdateMaintGrp(name string, description string, maintMaintGrpAttr models.MaintGrpAttributes) (*models.MaintGrp, error) {
	rn := fmt.Sprintf("maintgrp-%s", name)
	parentDn := fmt.Sprintf("uni/fabric")
	maintMaintGrp := models.NewMaintGrp(rn, parentDn, description, maintMaintGrpAttr)
	maintMaintGrp.Status = "modified"
	err := sm.Save(maintMaintGrp)
	return maintMaintGrp, err

}

func (sm *ServiceManager) ListMaintGrp() ([]*models.MaintGrp, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/maintMaintGrp.json", baseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MaintGrpListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationmaintRsMgrppFromMaintGrp(parentDn, tnMaintMaintPName string) error {
	dn := fmt.Sprintf("%s/rsmgrpp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMaintMaintPName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "maintRsMgrpp", dn, tnMaintMaintPName))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cont)

	return nil
}
