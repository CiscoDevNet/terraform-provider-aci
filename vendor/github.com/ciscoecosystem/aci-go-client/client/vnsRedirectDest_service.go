package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateDestinationofredirectedtraffic(ip string, service_redirect_policy string, tenant string, description string, vnsRedirectDestattr models.DestinationofredirectedtrafficAttributes) (*models.Destinationofredirectedtraffic, error) {
	rn := fmt.Sprintf("RedirectDest_ip-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", tenant, service_redirect_policy)
	vnsRedirectDest := models.NewDestinationofredirectedtraffic(rn, parentDn, description, vnsRedirectDestattr)
	err := sm.Save(vnsRedirectDest)
	return vnsRedirectDest, err
}

func (sm *ServiceManager) ReadDestinationofredirectedtraffic(ip string, service_redirect_policy string, tenant string) (*models.Destinationofredirectedtraffic, error) {
	dn := fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s/RedirectDest_ip-[%s]", tenant, service_redirect_policy, ip)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsRedirectDest := models.DestinationofredirectedtrafficFromContainer(cont)
	return vnsRedirectDest, nil
}

func (sm *ServiceManager) DeleteDestinationofredirectedtraffic(ip string, service_redirect_policy string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s/RedirectDest_ip-[%s]", tenant, service_redirect_policy, ip)
	return sm.DeleteByDn(dn, models.VnsredirectdestClassName)
}

func (sm *ServiceManager) UpdateDestinationofredirectedtraffic(ip string, service_redirect_policy string, tenant string, description string, vnsRedirectDestattr models.DestinationofredirectedtrafficAttributes) (*models.Destinationofredirectedtraffic, error) {
	rn := fmt.Sprintf("RedirectDest_ip-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", tenant, service_redirect_policy)
	vnsRedirectDest := models.NewDestinationofredirectedtraffic(rn, parentDn, description, vnsRedirectDestattr)

	vnsRedirectDest.Status = "modified"
	err := sm.Save(vnsRedirectDest)
	return vnsRedirectDest, err

}

func (sm *ServiceManager) ListDestinationofredirectedtraffic(service_redirect_policy string, tenant string) ([]*models.Destinationofredirectedtraffic, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/svcCont/svcRedirectPol-%s/vnsRedirectDest.json", baseurlStr, tenant, service_redirect_policy)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.DestinationofredirectedtrafficListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsRedirectHealthGroupFromDestinationofredirectedtraffic(parentDn, tnVnsRedirectHealthGroupName string) error {
	dn := fmt.Sprintf("%s/rsRedirectHealthGroup", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsRedirectHealthGroup", dn, tnVnsRedirectHealthGroupName))

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

func (sm *ServiceManager) DeleteRelationvnsRsRedirectHealthGroupFromDestinationofredirectedtraffic(parentDn string) error {
	dn := fmt.Sprintf("%s/rsRedirectHealthGroup", parentDn)
	return sm.DeleteByDn(dn, "vnsRsRedirectHealthGroup")
}

func (sm *ServiceManager) ReadRelationvnsRsRedirectHealthGroupFromDestinationofredirectedtraffic(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsRedirectHealthGroup")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsRedirectHealthGroup")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
