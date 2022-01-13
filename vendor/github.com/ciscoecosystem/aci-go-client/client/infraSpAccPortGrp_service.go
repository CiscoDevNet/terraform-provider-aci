package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSpineAccessPortPolicyGroup(name string, description string, infraSpAccPortGrpattr models.SpineAccessPortPolicyGroupAttributes) (*models.SpineAccessPortPolicyGroup, error) {
	rn := fmt.Sprintf("infra/funcprof/spaccportgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraSpAccPortGrp := models.NewSpineAccessPortPolicyGroup(rn, parentDn, description, infraSpAccPortGrpattr)
	err := sm.Save(infraSpAccPortGrp)
	return infraSpAccPortGrp, err
}

func (sm *ServiceManager) ReadSpineAccessPortPolicyGroup(name string) (*models.SpineAccessPortPolicyGroup, error) {
	dn := fmt.Sprintf("uni/infra/funcprof/spaccportgrp-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpAccPortGrp := models.SpineAccessPortPolicyGroupFromContainer(cont)
	return infraSpAccPortGrp, nil
}

func (sm *ServiceManager) DeleteSpineAccessPortPolicyGroup(name string) error {
	dn := fmt.Sprintf("uni/infra/funcprof/spaccportgrp-%s", name)
	return sm.DeleteByDn(dn, models.InfraspaccportgrpClassName)
}

func (sm *ServiceManager) UpdateSpineAccessPortPolicyGroup(name string, description string, infraSpAccPortGrpattr models.SpineAccessPortPolicyGroupAttributes) (*models.SpineAccessPortPolicyGroup, error) {
	rn := fmt.Sprintf("infra/funcprof/spaccportgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraSpAccPortGrp := models.NewSpineAccessPortPolicyGroup(rn, parentDn, description, infraSpAccPortGrpattr)

	infraSpAccPortGrp.Status = "modified"
	err := sm.Save(infraSpAccPortGrp)
	return infraSpAccPortGrp, err

}

func (sm *ServiceManager) ListSpineAccessPortPolicyGroup() ([]*models.SpineAccessPortPolicyGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infraSpAccPortGrp.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.SpineAccessPortPolicyGroupListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsHIfPolFromSpineAccessPortPolicyGroup(parentDn, tnFabricHIfPolName string) error {
	dn := fmt.Sprintf("%s/rshIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFabricHIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsHIfPol", dn, tnFabricHIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsHIfPolFromSpineAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsHIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsHIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsCdpIfPolFromSpineAccessPortPolicyGroup(parentDn, tnCdpIfPolName string) error {
	dn := fmt.Sprintf("%s/rscdpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnCdpIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsCdpIfPol", dn, tnCdpIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsCdpIfPolFromSpineAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsCdpIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsCdpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsCoppIfPolFromSpineAccessPortPolicyGroup(parentDn, tnCoppIfPolName string) error {
	dn := fmt.Sprintf("%s/rscoppIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnCoppIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsCoppIfPol", dn, tnCoppIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsCoppIfPolFromSpineAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsCoppIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsCoppIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsAttEntPFromSpineAccessPortPolicyGroup(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsattEntP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsAttEntP", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationinfraRsAttEntPFromSpineAccessPortPolicyGroup(parentDn string) error {
	dn := fmt.Sprintf("%s/rsattEntP", parentDn)
	return sm.DeleteByDn(dn, "infraRsAttEntP")
}

func (sm *ServiceManager) ReadRelationinfraRsAttEntPFromSpineAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsAttEntP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsAttEntP")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsMacsecIfPolFromSpineAccessPortPolicyGroup(parentDn, tnMacsecIfPolName string) error {
	dn := fmt.Sprintf("%s/rsmacsecIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMacsecIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsMacsecIfPol", dn, tnMacsecIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsMacsecIfPolFromSpineAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsMacsecIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsMacsecIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
