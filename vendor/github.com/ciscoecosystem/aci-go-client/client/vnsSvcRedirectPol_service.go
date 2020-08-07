package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateServiceRedirectPolicy(name string, tenant string, description string, vnsSvcRedirectPolattr models.ServiceRedirectPolicyAttributes) (*models.ServiceRedirectPolicy, error) {
	rn := fmt.Sprintf("svcCont/svcRedirectPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vnsSvcRedirectPol := models.NewServiceRedirectPolicy(rn, parentDn, description, vnsSvcRedirectPolattr)
	err := sm.Save(vnsSvcRedirectPol)
	return vnsSvcRedirectPol, err
}

func (sm *ServiceManager) ReadServiceRedirectPolicy(name string, tenant string) (*models.ServiceRedirectPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsSvcRedirectPol := models.ServiceRedirectPolicyFromContainer(cont)
	return vnsSvcRedirectPol, nil
}

func (sm *ServiceManager) DeleteServiceRedirectPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/svcCont/svcRedirectPol-%s", tenant, name)
	return sm.DeleteByDn(dn, models.VnssvcredirectpolClassName)
}

func (sm *ServiceManager) UpdateServiceRedirectPolicy(name string, tenant string, description string, vnsSvcRedirectPolattr models.ServiceRedirectPolicyAttributes) (*models.ServiceRedirectPolicy, error) {
	rn := fmt.Sprintf("svcCont/svcRedirectPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vnsSvcRedirectPol := models.NewServiceRedirectPolicy(rn, parentDn, description, vnsSvcRedirectPolattr)

	vnsSvcRedirectPol.Status = "modified"
	err := sm.Save(vnsSvcRedirectPol)
	return vnsSvcRedirectPol, err

}

func (sm *ServiceManager) ListServiceRedirectPolicy(tenant string) ([]*models.ServiceRedirectPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/vnsSvcRedirectPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ServiceRedirectPolicyListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsIPSLAMonitoringPolFromServiceRedirectPolicy(parentDn, tnFvIPSLAMonitoringPolName string) error {
	dn := fmt.Sprintf("%s/rsIPSLAMonitoringPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vnsRsIPSLAMonitoringPol", dn, tnFvIPSLAMonitoringPolName))

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

func (sm *ServiceManager) DeleteRelationvnsRsIPSLAMonitoringPolFromServiceRedirectPolicy(parentDn string) error {
	dn := fmt.Sprintf("%s/rsIPSLAMonitoringPol", parentDn)
	return sm.DeleteByDn(dn, "vnsRsIPSLAMonitoringPol")
}

func (sm *ServiceManager) ReadRelationvnsRsIPSLAMonitoringPolFromServiceRedirectPolicy(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vnsRsIPSLAMonitoringPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vnsRsIPSLAMonitoringPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
