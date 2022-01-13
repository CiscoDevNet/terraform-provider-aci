package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLogicalInterfaceContext(connNameOrLbl string, logicalDeviceContextNodeNameOrLbl string, logicalDeviceContextGraphNameOrLbl string, logical_device_context_ctrctNameOrLbl string, tenant string, description string, vnsLIfCtxattr models.LogicalInterfaceContextAttributes) (*models.LogicalInterfaceContext, error) {
	rn := fmt.Sprintf("lIfCtx-c-%s", connNameOrLbl)
	parentDn := fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-%s-n-%s", tenant, logical_device_context_ctrctNameOrLbl, logicalDeviceContextGraphNameOrLbl, logicalDeviceContextNodeNameOrLbl)
	vnsLIfCtx := models.NewLogicalInterfaceContext(rn, parentDn, description, vnsLIfCtxattr)
	err := sm.Save(vnsLIfCtx)
	return vnsLIfCtx, err
}

func (sm *ServiceManager) ReadLogicalInterfaceContext(connNameOrLbl string, logicalDeviceContextNodeNameOrLbl string, logicalDeviceContextGraphNameOrLbl string, logical_device_context_ctrctNameOrLbl string, tenant string) (*models.LogicalInterfaceContext, error) {
	dn := fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-%s-n-%s/lIfCtx-c-%s", tenant, logical_device_context_ctrctNameOrLbl, logicalDeviceContextGraphNameOrLbl, logicalDeviceContextNodeNameOrLbl, connNameOrLbl)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsLIfCtx := models.LogicalInterfaceContextFromContainer(cont)
	return vnsLIfCtx, nil
}

func (sm *ServiceManager) DeleteLogicalInterfaceContext(connNameOrLbl string, logicalDeviceContextNodeNameOrLbl string, logicalDeviceContextGraphNameOrLbl string, logical_device_context_ctrctNameOrLbl string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-%s-n-%s/lIfCtx-c-%s", tenant, logical_device_context_ctrctNameOrLbl, logicalDeviceContextGraphNameOrLbl, logicalDeviceContextNodeNameOrLbl, connNameOrLbl)
	return sm.DeleteByDn(dn, models.VnslifctxClassName)
}

func (sm *ServiceManager) UpdateLogicalInterfaceContext(connNameOrLbl string, logicalDeviceContextNodeNameOrLbl string, logicalDeviceContextGraphNameOrLbl string, logical_device_context_ctrctNameOrLbl string, tenant string, description string, vnsLIfCtxattr models.LogicalInterfaceContextAttributes) (*models.LogicalInterfaceContext, error) {
	rn := fmt.Sprintf("lIfCtx-c-%s", connNameOrLbl)
	parentDn := fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-%s-n-%s", tenant, logical_device_context_ctrctNameOrLbl, logicalDeviceContextGraphNameOrLbl, logicalDeviceContextNodeNameOrLbl)
	vnsLIfCtx := models.NewLogicalInterfaceContext(rn, parentDn, description, vnsLIfCtxattr)

	vnsLIfCtx.Status = "modified"
	err := sm.Save(vnsLIfCtx)
	return vnsLIfCtx, err

}

func (sm *ServiceManager) ListLogicalInterfaceContext(logicalDeviceContextNodeNameOrLbl string, logicalDeviceContextGraphNameOrLbl string, logical_device_context_ctrctNameOrLbl string, tenant string) ([]*models.LogicalInterfaceContext, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ldevCtx-c-%s-g-%s-n-%s/vnsLIfCtx.json", baseurlStr, tenant, logical_device_context_ctrctNameOrLbl, logicalDeviceContextGraphNameOrLbl, logicalDeviceContextNodeNameOrLbl)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LogicalInterfaceContextListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsLIfCtxToCustQosPolFromLogicalInterfaceContext(parentDn, tnQosCustomPolName string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToCustQosPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosCustomPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsLIfCtxToCustQosPol", dn, tnQosCustomPolName))

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

func (sm *ServiceManager) ReadRelationvnsRsLIfCtxToCustQosPolFromLogicalInterfaceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLIfCtxToCustQosPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLIfCtxToCustQosPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsLIfCtxToSvcEPgPolFromLogicalInterfaceContext(parentDn, tnVnsSvcEPgPolName string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToSvcEPgPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsLIfCtxToSvcEPgPol", dn, tnVnsSvcEPgPolName))

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

func (sm *ServiceManager) DeleteRelationvnsRsLIfCtxToSvcEPgPolFromLogicalInterfaceContext(parentDn string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToSvcEPgPol", parentDn)
	return sm.DeleteByDn(dn, "vnsRsLIfCtxToSvcEPgPol")
}

func (sm *ServiceManager) ReadRelationvnsRsLIfCtxToSvcEPgPolFromLogicalInterfaceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLIfCtxToSvcEPgPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLIfCtxToSvcEPgPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsLIfCtxToSvcRedirectPolFromLogicalInterfaceContext(parentDn, tnVnsSvcRedirectPolName string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToSvcRedirectPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsLIfCtxToSvcRedirectPol", dn, tnVnsSvcRedirectPolName))

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

func (sm *ServiceManager) DeleteRelationvnsRsLIfCtxToSvcRedirectPolFromLogicalInterfaceContext(parentDn string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToSvcRedirectPol", parentDn)
	return sm.DeleteByDn(dn, "vnsRsLIfCtxToSvcRedirectPol")
}

func (sm *ServiceManager) ReadRelationvnsRsLIfCtxToSvcRedirectPolFromLogicalInterfaceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLIfCtxToSvcRedirectPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLIfCtxToSvcRedirectPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsLIfCtxToLIfFromLogicalInterfaceContext(parentDn, tnVnsALDevLIfName string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToLIf", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsLIfCtxToLIf", dn, tnVnsALDevLIfName))

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

func (sm *ServiceManager) DeleteRelationvnsRsLIfCtxToLIfFromLogicalInterfaceContext(parentDn string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToLIf", parentDn)
	return sm.DeleteByDn(dn, "vnsRsLIfCtxToLIf")
}

func (sm *ServiceManager) ReadRelationvnsRsLIfCtxToLIfFromLogicalInterfaceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLIfCtxToLIf")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLIfCtxToLIf")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsLIfCtxToOutDefFromLogicalInterfaceContext(parentDn, tnL3extOutDefName string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToOutDef", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s"
								
			}
		}
	}`, "vnsRsLIfCtxToOutDef", dn, tnL3extOutDefName))

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

func (sm *ServiceManager) ReadRelationvnsRsLIfCtxToOutDefFromLogicalInterfaceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLIfCtxToOutDef")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLIfCtxToOutDef")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsLIfCtxToInstPFromLogicalInterfaceContext(parentDn, tnFvEPgName string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToInstP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s"
								
			}
		}
	}`, "vnsRsLIfCtxToInstP", dn, tnFvEPgName))

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

func (sm *ServiceManager) DeleteRelationvnsRsLIfCtxToInstPFromLogicalInterfaceContext(parentDn string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToInstP", parentDn)
	return sm.DeleteByDn(dn, "vnsRsLIfCtxToInstP")
}

func (sm *ServiceManager) ReadRelationvnsRsLIfCtxToInstPFromLogicalInterfaceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLIfCtxToInstP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLIfCtxToInstP")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsLIfCtxToBDFromLogicalInterfaceContext(parentDn, tnFvBDName string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToBD", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsLIfCtxToBD", dn, tnFvBDName))

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

func (sm *ServiceManager) DeleteRelationvnsRsLIfCtxToBDFromLogicalInterfaceContext(parentDn string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToBD", parentDn)
	return sm.DeleteByDn(dn, "vnsRsLIfCtxToBD")
}

func (sm *ServiceManager) ReadRelationvnsRsLIfCtxToBDFromLogicalInterfaceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLIfCtxToBD")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLIfCtxToBD")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsLIfCtxToOutFromLogicalInterfaceContext(parentDn, tnL3extOutName string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToOut", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsLIfCtxToOut", dn, tnL3extOutName))

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

func (sm *ServiceManager) DeleteRelationvnsRsLIfCtxToOutFromLogicalInterfaceContext(parentDn string) error {
	dn := fmt.Sprintf("%s/rsLIfCtxToOut", parentDn)
	return sm.DeleteByDn(dn, "vnsRsLIfCtxToOut")
}

func (sm *ServiceManager) ReadRelationvnsRsLIfCtxToOutFromLogicalInterfaceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLIfCtxToOut")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLIfCtxToOut")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
