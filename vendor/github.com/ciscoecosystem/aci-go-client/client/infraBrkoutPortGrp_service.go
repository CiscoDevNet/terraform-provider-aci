package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLeafBreakoutPortGroup(name string, description string, infraBrkoutPortGrpattr models.LeafBreakoutPortGroupAttributes) (*models.LeafBreakoutPortGroup, error) {
	rn := fmt.Sprintf("infra/funcprof/brkoutportgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraBrkoutPortGrp := models.NewLeafBreakoutPortGroup(rn, parentDn, description, infraBrkoutPortGrpattr)
	err := sm.Save(infraBrkoutPortGrp)
	return infraBrkoutPortGrp, err
}

func (sm *ServiceManager) ReadLeafBreakoutPortGroup(name string) (*models.LeafBreakoutPortGroup, error) {
	dn := fmt.Sprintf("uni/infra/funcprof/brkoutportgrp-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraBrkoutPortGrp := models.LeafBreakoutPortGroupFromContainer(cont)
	return infraBrkoutPortGrp, nil
}

func (sm *ServiceManager) DeleteLeafBreakoutPortGroup(name string) error {
	dn := fmt.Sprintf("uni/infra/funcprof/brkoutportgrp-%s", name)
	return sm.DeleteByDn(dn, models.InfrabrkoutportgrpClassName)
}

func (sm *ServiceManager) UpdateLeafBreakoutPortGroup(name string, description string, infraBrkoutPortGrpattr models.LeafBreakoutPortGroupAttributes) (*models.LeafBreakoutPortGroup, error) {
	rn := fmt.Sprintf("infra/funcprof/brkoutportgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraBrkoutPortGrp := models.NewLeafBreakoutPortGroup(rn, parentDn, description, infraBrkoutPortGrpattr)

	infraBrkoutPortGrp.Status = "modified"
	err := sm.Save(infraBrkoutPortGrp)
	return infraBrkoutPortGrp, err

}

func (sm *ServiceManager) ListLeafBreakoutPortGroup() ([]*models.LeafBreakoutPortGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infraBrkoutPortGrp.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LeafBreakoutPortGroupListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsMonBrkoutInfraPolFromLeafBreakoutPortGroup(parentDn, tnMonInfraPolName string) error {
	dn := fmt.Sprintf("%s/rsmonBrkoutInfraPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMonInfraPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsMonBrkoutInfraPol", dn, tnMonInfraPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsMonBrkoutInfraPolFromLeafBreakoutPortGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsMonBrkoutInfraPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsMonBrkoutInfraPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
