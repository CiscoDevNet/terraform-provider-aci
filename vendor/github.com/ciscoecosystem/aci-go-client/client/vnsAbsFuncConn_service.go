package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFunctionConnector(name string, function_node string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsFuncConnattr models.FunctionConnectorAttributes) (*models.FunctionConnector, error) {
	rn := fmt.Sprintf("AbsFConn-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s", tenant, l4_l7_service_graph_template, function_node)
	vnsAbsFuncConn := models.NewFunctionConnector(rn, parentDn, description, vnsAbsFuncConnattr)
	err := sm.Save(vnsAbsFuncConn)
	return vnsAbsFuncConn, err
}

func (sm *ServiceManager) ReadFunctionConnector(name string, function_node string, l4_l7_service_graph_template string, tenant string) (*models.FunctionConnector, error) {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s/AbsFConn-%s", tenant, l4_l7_service_graph_template, function_node, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsFuncConn := models.FunctionConnectorFromContainer(cont)
	return vnsAbsFuncConn, nil
}

func (sm *ServiceManager) DeleteFunctionConnector(name string, function_node string, l4_l7_service_graph_template string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s/AbsFConn-%s", tenant, l4_l7_service_graph_template, function_node, name)
	return sm.DeleteByDn(dn, models.VnsabsfuncconnClassName)
}

func (sm *ServiceManager) UpdateFunctionConnector(name string, function_node string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsFuncConnattr models.FunctionConnectorAttributes) (*models.FunctionConnector, error) {
	rn := fmt.Sprintf("AbsFConn-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s", tenant, l4_l7_service_graph_template, function_node)
	vnsAbsFuncConn := models.NewFunctionConnector(rn, parentDn, description, vnsAbsFuncConnattr)

	vnsAbsFuncConn.Status = "modified"
	err := sm.Save(vnsAbsFuncConn)
	return vnsAbsFuncConn, err

}

func (sm *ServiceManager) ListFunctionConnector(function_node string, l4_l7_service_graph_template string, tenant string) ([]*models.FunctionConnector, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/AbsGraph-%s/AbsNode-%s/vnsAbsFuncConn.json", baseurlStr, tenant, l4_l7_service_graph_template, function_node)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FunctionConnectorListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsConnToCtxTermFromFunctionConnector(parentDn, tnVnsATermName string) error {
	dn := fmt.Sprintf("%s/rsConnToCtxTerm", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVnsATermName": "%s"
								
			}
		}
	}`, "vnsRsConnToCtxTerm", dn, tnVnsATermName))

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

func (sm *ServiceManager) DeleteRelationvnsRsConnToCtxTermFromFunctionConnector(parentDn string) error {
	dn := fmt.Sprintf("%s/rsConnToCtxTerm", parentDn)
	return sm.DeleteByDn(dn, "vnsRsConnToCtxTerm")
}

func (sm *ServiceManager) ReadRelationvnsRsConnToCtxTermFromFunctionConnector(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsConnToCtxTerm")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsConnToCtxTerm")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnVnsATermName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsConnToFltFromFunctionConnector(parentDn, tnVzAFilterableName string) error {
	dn := fmt.Sprintf("%s/rsConnToFlt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVzAFilterableName": "%s"
								
			}
		}
	}`, "vnsRsConnToFlt", dn, tnVzAFilterableName))

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

func (sm *ServiceManager) DeleteRelationvnsRsConnToFltFromFunctionConnector(parentDn string) error {
	dn := fmt.Sprintf("%s/rsConnToFlt", parentDn)
	return sm.DeleteByDn(dn, "vnsRsConnToFlt")
}

func (sm *ServiceManager) ReadRelationvnsRsConnToFltFromFunctionConnector(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsConnToFlt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsConnToFlt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnVzAFilterableName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsMConnAttFromFunctionConnector(parentDn, tnVnsMConnName string) error {
	dn := fmt.Sprintf("%s/rsMConnAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVnsMConnName": "%s"
								
			}
		}
	}`, "vnsRsMConnAtt", dn, tnVnsMConnName))

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

func (sm *ServiceManager) DeleteRelationvnsRsMConnAttFromFunctionConnector(parentDn string) error {
	dn := fmt.Sprintf("%s/rsMConnAtt", parentDn)
	return sm.DeleteByDn(dn, "vnsRsMConnAtt")
}

func (sm *ServiceManager) ReadRelationvnsRsMConnAttFromFunctionConnector(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsMConnAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsMConnAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnVnsMConnName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsConnToAConnFromFunctionConnector(parentDn, tnVnsAConnName string) error {
	dn := fmt.Sprintf("%s/rsConnToAConn", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVnsAConnName": "%s"
								
			}
		}
	}`, "vnsRsConnToAConn", dn, tnVnsAConnName))

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

func (sm *ServiceManager) DeleteRelationvnsRsConnToAConnFromFunctionConnector(parentDn string) error {
	dn := fmt.Sprintf("%s/rsConnToAConn", parentDn)
	return sm.DeleteByDn(dn, "vnsRsConnToAConn")
}

func (sm *ServiceManager) ReadRelationvnsRsConnToAConnFromFunctionConnector(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsConnToAConn")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsConnToAConn")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnVnsAConnName")
		return dat, err
	} else {
		return nil, err
	}

}
