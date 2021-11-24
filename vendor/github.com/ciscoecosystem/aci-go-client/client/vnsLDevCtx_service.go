package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLogicalDeviceContext(nodeNameOrLbl string, graphNameOrLbl string, ctrctNameOrLbl string, tenant string, description string, vnsLDevCtxattr models.LogicalDeviceContextAttributes) (*models.LogicalDeviceContext, error) {
	rn := fmt.Sprintf("ldevCtx-c-%s-g-%s-n-%s", ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vnsLDevCtx := models.NewLogicalDeviceContext(rn, parentDn, description, vnsLDevCtxattr)
	err := sm.Save(vnsLDevCtx)
	return vnsLDevCtx, err
}

func (sm *ServiceManager) ReadLogicalDeviceContext(nodeNameOrLbl string, graphNameOrLbl string, ctrctNameOrLbl string, tenant string) (*models.LogicalDeviceContext, error) {
	dn := fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-%s-n-%s", tenant, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsLDevCtx := models.LogicalDeviceContextFromContainer(cont)
	return vnsLDevCtx, nil
}

func (sm *ServiceManager) DeleteLogicalDeviceContext(nodeNameOrLbl string, graphNameOrLbl string, ctrctNameOrLbl string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-%s-n-%s", tenant, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	return sm.DeleteByDn(dn, models.VnsldevctxClassName)
}

func (sm *ServiceManager) UpdateLogicalDeviceContext(nodeNameOrLbl string, graphNameOrLbl string, ctrctNameOrLbl string, tenant string, description string, vnsLDevCtxattr models.LogicalDeviceContextAttributes) (*models.LogicalDeviceContext, error) {
	rn := fmt.Sprintf("ldevCtx-c-%s-g-%s-n-%s", ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vnsLDevCtx := models.NewLogicalDeviceContext(rn, parentDn, description, vnsLDevCtxattr)

	vnsLDevCtx.Status = "modified"
	err := sm.Save(vnsLDevCtx)
	return vnsLDevCtx, err

}

func (sm *ServiceManager) ListLogicalDeviceContext(tenant string) ([]*models.LogicalDeviceContext, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/vnsLDevCtx.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LogicalDeviceContextListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsLDevCtxToLDevFromLogicalDeviceContext(parentDn, tnVnsALDevIfName string) error {
	dn := fmt.Sprintf("%s/rsLDevCtxToLDev", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsLDevCtxToLDev", dn, tnVnsALDevIfName))

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

func (sm *ServiceManager) DeleteRelationvnsRsLDevCtxToLDevFromLogicalDeviceContext(parentDn string) error {
	dn := fmt.Sprintf("%s/rsLDevCtxToLDev", parentDn)
	return sm.DeleteByDn(dn, "vnsRsLDevCtxToLDev")
}

func (sm *ServiceManager) ReadRelationvnsRsLDevCtxToLDevFromLogicalDeviceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLDevCtxToLDev")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLDevCtxToLDev")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvnsRsLDevCtxToRtrCfgFromLogicalDeviceContext(parentDn, tnVnsRtrCfgName string) error {
	dn := fmt.Sprintf("%s/rsLDevCtxToRtrCfg", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVnsRtrCfgName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsLDevCtxToRtrCfg", dn, tnVnsRtrCfgName))

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

func (sm *ServiceManager) DeleteRelationvnsRsLDevCtxToRtrCfgFromLogicalDeviceContext(parentDn string) error {
	dn := fmt.Sprintf("%s/rsLDevCtxToRtrCfg", parentDn)
	return sm.DeleteByDn(dn, "vnsRsLDevCtxToRtrCfg")
}

func (sm *ServiceManager) ReadRelationvnsRsLDevCtxToRtrCfgFromLogicalDeviceContext(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsLDevCtxToRtrCfg")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsLDevCtxToRtrCfg")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
