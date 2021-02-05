package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFunctionNode(name string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsNodeattr models.FunctionNodeAttributes) (*models.FunctionNode, error) {
	rn := fmt.Sprintf("AbsNode-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, l4_l7_service_graph_template)
	vnsAbsNode := models.NewFunctionNode(rn, parentDn, description, vnsAbsNodeattr)
	err := sm.Save(vnsAbsNode)
	return vnsAbsNode, err
}

func (sm *ServiceManager) ReadFunctionNode(name string, l4_l7_service_graph_template string, tenant string) (*models.FunctionNode, error) {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s", tenant, l4_l7_service_graph_template, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsNode := models.FunctionNodeFromContainer(cont)
	return vnsAbsNode, nil
}

func (sm *ServiceManager) DeleteFunctionNode(name string, l4_l7_service_graph_template string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s", tenant, l4_l7_service_graph_template, name)
	return sm.DeleteByDn(dn, models.VnsabsnodeClassName)
}

func (sm *ServiceManager) UpdateFunctionNode(name string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsNodeattr models.FunctionNodeAttributes) (*models.FunctionNode, error) {
	rn := fmt.Sprintf("AbsNode-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, l4_l7_service_graph_template)
	vnsAbsNode := models.NewFunctionNode(rn, parentDn, description, vnsAbsNodeattr)

	vnsAbsNode.Status = "modified"
	err := sm.Save(vnsAbsNode)
	return vnsAbsNode, err

}

func (sm *ServiceManager) ListFunctionNode(l4_l7_service_graph_template string, tenant string) ([]*models.FunctionNode, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/AbsGraph-%s/vnsAbsNode.json", baseurlStr, tenant, l4_l7_service_graph_template)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FunctionNodeListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsNodeToAbsFuncProfFromFunctionNode(parentDn, tnVnsAbsFuncProfName string) error {
	dn := fmt.Sprintf("%s/rsNodeToAbsFuncProf", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsNodeToAbsFuncProf", dn, tnVnsAbsFuncProfName))

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

func (sm *ServiceManager) DeleteRelationvnsRsNodeToAbsFuncProfFromFunctionNode(parentDn string) error {
	dn := fmt.Sprintf("%s/rsNodeToAbsFuncProf", parentDn)
	return sm.DeleteByDn(dn, "vnsRsNodeToAbsFuncProf")
}

func (sm *ServiceManager) ReadRelationvnsRsNodeToAbsFuncProfFromFunctionNode(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsNodeToAbsFuncProf")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsNodeToAbsFuncProf")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsNodeToLDevFromFunctionNode(parentDn, tnVnsALDevIfName string) error {
	dn := fmt.Sprintf("%s/rsNodeToLDev", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsNodeToLDev", dn, tnVnsALDevIfName))

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

func (sm *ServiceManager) DeleteRelationvnsRsNodeToLDevFromFunctionNode(parentDn string) error {
	dn := fmt.Sprintf("%s/rsNodeToLDev", parentDn)
	return sm.DeleteByDn(dn, "vnsRsNodeToLDev")
}

func (sm *ServiceManager) ReadRelationvnsRsNodeToLDevFromFunctionNode(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsNodeToLDev")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsNodeToLDev")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsNodeToMFuncFromFunctionNode(parentDn, tnVnsMFuncName string) error {
	dn := fmt.Sprintf("%s/rsNodeToMFunc", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsNodeToMFunc", dn, tnVnsMFuncName))

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

func (sm *ServiceManager) DeleteRelationvnsRsNodeToMFuncFromFunctionNode(parentDn string) error {
	dn := fmt.Sprintf("%s/rsNodeToMFunc", parentDn)
	return sm.DeleteByDn(dn, "vnsRsNodeToMFunc")
}

func (sm *ServiceManager) ReadRelationvnsRsNodeToMFuncFromFunctionNode(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsNodeToMFunc")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsNodeToMFunc")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsDefaultScopeToTermFromFunctionNode(parentDn, tnVnsATermName string) error {
	dn := fmt.Sprintf("%s/rsdefaultScopeToTerm", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsDefaultScopeToTerm", dn, tnVnsATermName))

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

func (sm *ServiceManager) DeleteRelationvnsRsDefaultScopeToTermFromFunctionNode(parentDn string) error {
	dn := fmt.Sprintf("%s/rsdefaultScopeToTerm", parentDn)
	return sm.DeleteByDn(dn, "vnsRsDefaultScopeToTerm")
}

func (sm *ServiceManager) ReadRelationvnsRsDefaultScopeToTermFromFunctionNode(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsDefaultScopeToTerm")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsDefaultScopeToTerm")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsNodeToCloudLDevFromFunctionNode(parentDn, tnCloudALDevName string) error {
	dn := fmt.Sprintf("%s/rsNodeToCloudLDev", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsNodeToCloudLDev", dn, tnCloudALDevName))

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

func (sm *ServiceManager) DeleteRelationvnsRsNodeToCloudLDevFromFunctionNode(parentDn string) error {
	dn := fmt.Sprintf("%s/rsNodeToCloudLDev", parentDn)
	return sm.DeleteByDn(dn, "vnsRsNodeToCloudLDev")
}

func (sm *ServiceManager) ReadRelationvnsRsNodeToCloudLDevFromFunctionNode(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsNodeToCloudLDev")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsNodeToCloudLDev")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
