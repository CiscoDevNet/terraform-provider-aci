package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateLeafAccessPortPolicyGroup(name string, description string, infraAccPortGrpattr models.LeafAccessPortPolicyGroupAttributes) (*models.LeafAccessPortPolicyGroup, error) {
	rn := fmt.Sprintf("infra/funcprof/accportgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraAccPortGrp := models.NewLeafAccessPortPolicyGroup(rn, parentDn, description, infraAccPortGrpattr)
	err := sm.Save(infraAccPortGrp)
	return infraAccPortGrp, err
}

func (sm *ServiceManager) ReadLeafAccessPortPolicyGroup(name string) (*models.LeafAccessPortPolicyGroup, error) {
	dn := fmt.Sprintf("uni/infra/funcprof/accportgrp-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraAccPortGrp := models.LeafAccessPortPolicyGroupFromContainer(cont)
	return infraAccPortGrp, nil
}

func (sm *ServiceManager) DeleteLeafAccessPortPolicyGroup(name string) error {
	dn := fmt.Sprintf("uni/infra/funcprof/accportgrp-%s", name)
	return sm.DeleteByDn(dn, models.InfraaccportgrpClassName)
}

func (sm *ServiceManager) UpdateLeafAccessPortPolicyGroup(name string, description string, infraAccPortGrpattr models.LeafAccessPortPolicyGroupAttributes) (*models.LeafAccessPortPolicyGroup, error) {
	rn := fmt.Sprintf("infra/funcprof/accportgrp-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraAccPortGrp := models.NewLeafAccessPortPolicyGroup(rn, parentDn, description, infraAccPortGrpattr)

	infraAccPortGrp.Status = "modified"
	err := sm.Save(infraAccPortGrp)
	return infraAccPortGrp, err

}

func (sm *ServiceManager) ListLeafAccessPortPolicyGroup() ([]*models.LeafAccessPortPolicyGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infraAccPortGrp.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LeafAccessPortPolicyGroupListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(parentDn, tnSpanVSrcGrpName string) error {
	dn := fmt.Sprintf("%s/rsspanVSrcGrp-%s", parentDn, tnSpanVSrcGrpName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "infraRsSpanVSrcGrp", dn))

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

func (sm *ServiceManager) DeleteRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(parentDn, tnSpanVSrcGrpName string) error {
	dn := fmt.Sprintf("%s/rsspanVSrcGrp-%s", parentDn, tnSpanVSrcGrpName)
	return sm.DeleteByDn(dn, "infraRsSpanVSrcGrp")
}

func (sm *ServiceManager) ReadRelationinfraRsSpanVSrcGrpFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsSpanVSrcGrp")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsSpanVSrcGrp")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationinfraRsStormctrlIfPolFromLeafAccessPortPolicyGroup(parentDn, tnStormctrlIfPolName string) error {
	dn := fmt.Sprintf("%s/rsstormctrlIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnStormctrlIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsStormctrlIfPol", dn, tnStormctrlIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsStormctrlIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsStormctrlIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsStormctrlIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(parentDn, tnPoeIfPolName string) error {
	dn := fmt.Sprintf("%s/rspoeIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnPoeIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsPoeIfPol", dn, tnPoeIfPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(parentDn string) error {
	dn := fmt.Sprintf("%s/rspoeIfPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsPoeIfPol")
}

func (sm *ServiceManager) ReadRelationinfraRsPoeIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsPoeIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsPoeIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsLldpIfPolFromLeafAccessPortPolicyGroup(parentDn, tnLldpIfPolName string) error {
	dn := fmt.Sprintf("%s/rslldpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnLldpIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsLldpIfPol", dn, tnLldpIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsLldpIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsLldpIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsLldpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsMacsecIfPolFromLeafAccessPortPolicyGroup(parentDn, tnMacsecIfPolName string) error {
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) ReadRelationinfraRsMacsecIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationinfraRsQosDppIfPolFromLeafAccessPortPolicyGroup(parentDn, tnQosDppPolName string) error {
	dn := fmt.Sprintf("%s/rsqosDppIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosDppPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsQosDppIfPol", dn, tnQosDppPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsQosDppIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsQosDppIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsQosDppIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsHIfPolFromLeafAccessPortPolicyGroup(parentDn, tnFabricHIfPolName string) error {
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) ReadRelationinfraRsHIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(parentDn, tnNetflowMonitorPolName, fltType string) error {
	dn := fmt.Sprintf("%s/rsnetflowMonitorPol-[%s]-%s", parentDn, tnNetflowMonitorPolName, fltType)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "infraRsNetflowMonitorPol", dn))

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

	return CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)
}

func (sm *ServiceManager) DeleteRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(parentDn, tnNetflowMonitorPolName, fltType string) error {
	dn := fmt.Sprintf("%s/rsnetflowMonitorPol-[%s]-%s", parentDn, tnNetflowMonitorPolName, fltType)
	return sm.DeleteByDn(dn, "infraRsNetflowMonitorPol")
}

func (sm *ServiceManager) ReadRelationinfraRsNetflowMonitorPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsNetflowMonitorPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsNetflowMonitorPol")

	st := make([]map[string]string, 0)

	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnNetflowMonitorPolName"] = models.G(contItem, "tDn")
		paramMap["fltType"] = models.G(contItem, "fltType")

		st = append(st, paramMap)

	}

	return st, err

}
func (sm *ServiceManager) CreateRelationinfraRsL2PortAuthPolFromLeafAccessPortPolicyGroup(parentDn, tnL2PortAuthPolName string) error {
	dn := fmt.Sprintf("%s/rsl2PortAuthPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnL2PortAuthPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsL2PortAuthPol", dn, tnL2PortAuthPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsL2PortAuthPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsL2PortAuthPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsL2PortAuthPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsMcpIfPolFromLeafAccessPortPolicyGroup(parentDn, tnMcpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsmcpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMcpIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsMcpIfPol", dn, tnMcpIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsMcpIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsMcpIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsMcpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsL2PortSecurityPolFromLeafAccessPortPolicyGroup(parentDn, tnL2PortSecurityPolName string) error {
	dn := fmt.Sprintf("%s/rsl2PortSecurityPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnL2PortSecurityPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsL2PortSecurityPol", dn, tnL2PortSecurityPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsL2PortSecurityPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsL2PortSecurityPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsL2PortSecurityPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsCoppIfPolFromLeafAccessPortPolicyGroup(parentDn, tnCoppIfPolName string) error {
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) ReadRelationinfraRsCoppIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(parentDn, tnSpanVDestGrpName string) error {
	dn := fmt.Sprintf("%s/rsspanVDestGrp-%s", parentDn, tnSpanVDestGrpName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "infraRsSpanVDestGrp", dn))

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

func (sm *ServiceManager) DeleteRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(parentDn, tnSpanVDestGrpName string) error {
	dn := fmt.Sprintf("%s/rsspanVDestGrp-%s", parentDn, tnSpanVDestGrpName)
	return sm.DeleteByDn(dn, "infraRsSpanVDestGrp")
}

func (sm *ServiceManager) ReadRelationinfraRsSpanVDestGrpFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsSpanVDestGrp")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsSpanVDestGrp")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationinfraRsDwdmIfPolFromLeafAccessPortPolicyGroup(parentDn, tnDwdmIfPolName string) error {
	dn := fmt.Sprintf("%s/rsdwdmIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnDwdmIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsDwdmIfPol", dn, tnDwdmIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsDwdmIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsDwdmIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsDwdmIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsQosPfcIfPolFromLeafAccessPortPolicyGroup(parentDn, tnQosPfcIfPolName string) error {
	dn := fmt.Sprintf("%s/rsqosPfcIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosPfcIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsQosPfcIfPol", dn, tnQosPfcIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsQosPfcIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsQosPfcIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsQosPfcIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsQosSdIfPolFromLeafAccessPortPolicyGroup(parentDn, tnQosSdIfPolName string) error {
	dn := fmt.Sprintf("%s/rsqosSdIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosSdIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsQosSdIfPol", dn, tnQosSdIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsQosSdIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsQosSdIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsQosSdIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsMonIfInfraPolFromLeafAccessPortPolicyGroup(parentDn, tnMonInfraPolName string) error {
	dn := fmt.Sprintf("%s/rsmonIfInfraPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMonInfraPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsMonIfInfraPol", dn, tnMonInfraPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsMonIfInfraPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsMonIfInfraPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsMonIfInfraPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsFcIfPolFromLeafAccessPortPolicyGroup(parentDn, tnFcIfPolName string) error {
	dn := fmt.Sprintf("%s/rsfcIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFcIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsFcIfPol", dn, tnFcIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsFcIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsFcIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsFcIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsQosIngressDppIfPolFromLeafAccessPortPolicyGroup(parentDn, tnQosDppPolName string) error {
	dn := fmt.Sprintf("%s/rsQosIngressDppIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosDppPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsQosIngressDppIfPol", dn, tnQosDppPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsQosIngressDppIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsQosIngressDppIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsQosIngressDppIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsCdpIfPolFromLeafAccessPortPolicyGroup(parentDn, tnCdpIfPolName string) error {
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) ReadRelationinfraRsCdpIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationinfraRsL2IfPolFromLeafAccessPortPolicyGroup(parentDn, tnL2IfPolName string) error {
	dn := fmt.Sprintf("%s/rsl2IfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnL2IfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsL2IfPol", dn, tnL2IfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsL2IfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsL2IfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsL2IfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsStpIfPolFromLeafAccessPortPolicyGroup(parentDn, tnStpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsstpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnStpIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsStpIfPol", dn, tnStpIfPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsStpIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsStpIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsStpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsQosEgressDppIfPolFromLeafAccessPortPolicyGroup(parentDn, tnQosDppPolName string) error {
	dn := fmt.Sprintf("%s/rsQosEgressDppIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosDppPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsQosEgressDppIfPol", dn, tnQosDppPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsQosEgressDppIfPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsQosEgressDppIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsQosEgressDppIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(parentDn, tDn string) error {
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) DeleteRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(parentDn string) error {
	dn := fmt.Sprintf("%s/rsattEntP", parentDn)
	return sm.DeleteByDn(dn, "infraRsAttEntP")
}

func (sm *ServiceManager) ReadRelationinfraRsAttEntPFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(parentDn, tnL2InstPolName string) error {
	dn := fmt.Sprintf("%s/rsl2InstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsL2InstPol", dn, tnL2InstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(parentDn string) error {
	dn := fmt.Sprintf("%s/rsl2InstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsL2InstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsL2InstPolFromLeafAccessPortPolicyGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsL2InstPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsL2InstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
