package client

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRouteContextScope(route_control_context string, route_control_profile string, tenant string, description string, nameAlias string, rtctrlScopeAttr models.RouteContextScopeAttributes) (*models.RouteContextScope, error) {
	rn := fmt.Sprintf(models.RnrtctrlScope)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlScope, tenant, route_control_profile, route_control_context)
	rtctrlScope := models.NewRouteContextScope(rn, parentDn, description, nameAlias, rtctrlScopeAttr)
	err := sm.Save(rtctrlScope)
	return rtctrlScope, err
}

func (sm *ServiceManager) ReadRouteContextScope(route_control_context string, route_control_profile string, tenant string) (*models.RouteContextScope, error) {
	dn := fmt.Sprintf(models.DnrtctrlScope, tenant, route_control_profile, route_control_context)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlScope := models.RouteContextScopeFromContainer(cont)
	return rtctrlScope, nil
}

func (sm *ServiceManager) DeleteRouteContextScope(route_control_context string, route_control_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlScope, tenant, route_control_profile, route_control_context)
	return sm.DeleteByDn(dn, models.RtctrlscopeClassName)
}

func (sm *ServiceManager) UpdateRouteContextScope(route_control_context string, route_control_profile string, tenant string, description string, nameAlias string, rtctrlScopeAttr models.RouteContextScopeAttributes) (*models.RouteContextScope, error) {
	rn := fmt.Sprintf(models.RnrtctrlScope)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlScope, tenant, route_control_profile, route_control_context)
	rtctrlScope := models.NewRouteContextScope(rn, parentDn, description, nameAlias, rtctrlScopeAttr)
	rtctrlScope.Status = "modified"
	err := sm.Save(rtctrlScope)
	return rtctrlScope, err
}

func (sm *ServiceManager) ListRouteContextScope(route_control_context string, route_control_profile string, tenant string) ([]*models.RouteContextScope, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/prof-%s/ctx-%s/rtctrlScope.json", models.BaseurlStr, tenant, route_control_profile, route_control_context)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RouteContextScopeListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationrtctrlRsScopeToAttrP(parentDn, annotation, tnRtctrlAttrPName string) error {
	dn := fmt.Sprintf("%s/rsScopeToAttrP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnRtctrlAttrPName": "%s"
			}
		}
	}`, "rtctrlRsScopeToAttrP", dn, annotation, tnRtctrlAttrPName))

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
	log.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationrtctrlRsScopeToAttrP(parentDn string) error {
	dn := fmt.Sprintf("%s/rsScopeToAttrP", parentDn)
	return sm.DeleteByDn(dn, "rtctrlRsScopeToAttrP")
}

func (sm *ServiceManager) ReadRelationrtctrlRsScopeToAttrP(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "rtctrlRsScopeToAttrP")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "rtctrlRsScopeToAttrP")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
