package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateLogicalInterfaceProfile(name string, logical_node_profile string, l3_outside string, tenant string, description string, l3extLIfPattr models.LogicalInterfaceProfileAttributes) (*models.LogicalInterfaceProfile, error) {
	rn := fmt.Sprintf("lifp-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", tenant, l3_outside, logical_node_profile)
	l3extLIfP := models.NewLogicalInterfaceProfile(rn, parentDn, description, l3extLIfPattr)
	err := sm.Save(l3extLIfP)
	return l3extLIfP, err
}

func (sm *ServiceManager) ReadLogicalInterfaceProfile(name string, logical_node_profile string, l3_outside string, tenant string) (*models.LogicalInterfaceProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extLIfP := models.LogicalInterfaceProfileFromContainer(cont)
	return l3extLIfP, nil
}

func (sm *ServiceManager) DeleteLogicalInterfaceProfile(name string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, name)
	return sm.DeleteByDn(dn, models.L3extlifpClassName)
}

func (sm *ServiceManager) UpdateLogicalInterfaceProfile(name string, logical_node_profile string, l3_outside string, tenant string, description string, l3extLIfPattr models.LogicalInterfaceProfileAttributes) (*models.LogicalInterfaceProfile, error) {
	rn := fmt.Sprintf("lifp-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", tenant, l3_outside, logical_node_profile)
	l3extLIfP := models.NewLogicalInterfaceProfile(rn, parentDn, description, l3extLIfPattr)

	l3extLIfP.Status = "modified"
	err := sm.Save(l3extLIfP)
	return l3extLIfP, err

}

func (sm *ServiceManager) ListLogicalInterfaceProfile(logical_node_profile string, l3_outside string, tenant string) ([]*models.LogicalInterfaceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/l3extLIfP.json", baseurlStr, tenant, l3_outside, logical_node_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LogicalInterfaceProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(parentDn, tnNetflowMonitorPolName, fltType string) error {
	dn := fmt.Sprintf("%s/rslIfPToNetflowMonitorPol-[%s]-%s", parentDn, tnNetflowMonitorPolName, fltType)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "l3extRsLIfPToNetflowMonitorPol", dn))

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

func (sm *ServiceManager) DeleteRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(parentDn, tnNetflowMonitorPolName, fltType string) error {
	dn := fmt.Sprintf("%s/rslIfPToNetflowMonitorPol-[%s]-%s", parentDn, tnNetflowMonitorPolName, fltType)
	return sm.DeleteByDn(dn, "l3extRsLIfPToNetflowMonitorPol")
}

func (sm *ServiceManager) ReadRelationl3extRsLIfPToNetflowMonitorPolFromLogicalInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsLIfPToNetflowMonitorPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsLIfPToNetflowMonitorPol")

	st := make([]map[string]string, 0)

	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnNetflowMonitorPolName"] = models.G(contItem, "tDn")
		paramMap["fltType"] = models.G(contItem, "fltType")

		st = append(st, paramMap)

	}

	return st, err

}
func (sm *ServiceManager) CreateRelationl3extRsPathL3OutAttFromLogicalInterfaceProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rspathL3OutAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "l3extRsPathL3OutAtt", dn))

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

func (sm *ServiceManager) DeleteRelationl3extRsPathL3OutAttFromLogicalInterfaceProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rspathL3OutAtt-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "l3extRsPathL3OutAtt")
}

func (sm *ServiceManager) ReadRelationl3extRsPathL3OutAttFromLogicalInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsPathL3OutAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsPathL3OutAtt")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(parentDn, tnQosDppPolName string) error {
	dn := fmt.Sprintf("%s/rsegressQosDppPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosDppPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsEgressQosDppPol", dn, tnQosDppPolName))

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

func (sm *ServiceManager) ReadRelationl3extRsEgressQosDppPolFromLogicalInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsEgressQosDppPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsEgressQosDppPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(parentDn, tnQosDppPolName string) error {
	dn := fmt.Sprintf("%s/rsingressQosDppPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosDppPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsIngressQosDppPol", dn, tnQosDppPolName))

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

func (sm *ServiceManager) ReadRelationl3extRsIngressQosDppPolFromLogicalInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsIngressQosDppPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsIngressQosDppPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(parentDn, tnQosCustomPolName string) error {
	dn := fmt.Sprintf("%s/rslIfPCustQosPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosCustomPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsLIfPCustQosPol", dn, tnQosCustomPolName))

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

func (sm *ServiceManager) ReadRelationl3extRsLIfPCustQosPolFromLogicalInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsLIfPCustQosPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsLIfPCustQosPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationl3extRsArpIfPolFromLogicalInterfaceProfile(parentDn, tnArpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsArpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnArpIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsArpIfPol", dn, tnArpIfPolName))

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

func (sm *ServiceManager) ReadRelationl3extRsArpIfPolFromLogicalInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsArpIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsArpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationl3extRsNdIfPolFromLogicalInterfaceProfile(parentDn, tnNdIfPolName string) error {
	dn := fmt.Sprintf("%s/rsNdIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnNdIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsNdIfPol", dn, tnNdIfPolName))

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

func (sm *ServiceManager) ReadRelationl3extRsNdIfPolFromLogicalInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsNdIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsNdIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
