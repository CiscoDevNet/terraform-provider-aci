package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudVpnGateway(name string, cloudContextProfile string, tenant string, description string, cloudRouterPattr models.CloudVpnGatewayAttributes) (*models.CloudVpnGateway, error) {
	rn := fmt.Sprintf("routerp-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s", tenant, cloudContextProfile)
	cloudRouterP := models.NewCloudVpnGateway(rn, parentDn, description, cloudRouterPattr)
	err := sm.Save(cloudRouterP)
	return cloudRouterP, err
}

func (sm *ServiceManager) ReadCloudVpnGateway(name string, cloudContextProfile string, tenant string) (*models.CloudVpnGateway, error) {
	dn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s/routerp-%s", tenant, cloudContextProfile, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudRouterP := models.CloudVpnGatewayFromContainer(cont)
	return cloudRouterP, nil
}

func (sm *ServiceManager) DeleteCloudVpnGateway(name string, cloudContextProfile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s/routerp-%s", tenant, cloudContextProfile, name)
	return sm.DeleteByDn(dn, models.CloudrouterpClassName)
}

func (sm *ServiceManager) UpdateCloudVpnGateway(name string, cloudContextProfile string, tenant string, description string, cloudRouterPattr models.CloudVpnGatewayAttributes) (*models.CloudVpnGateway, error) {
	rn := fmt.Sprintf("routerp-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s", tenant, cloudContextProfile)
	cloudRouterP := models.NewCloudVpnGateway(rn, parentDn, description, cloudRouterPattr)

	cloudRouterP.Status = "modified"
	err := sm.Save(cloudRouterP)
	return cloudRouterP, err

}

func (sm *ServiceManager) ListCloudVpnGateway(cloudContextProfile string, tenant string) ([]*models.CloudVpnGateway, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctxprofile-%s/cloudRouterP.json", baseurlStr, tenant, cloudContextProfile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudVpnGatewayListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationcloudRsToVpnGwPolFromCloudVpnGateway(parentDn, tnCloudVpnGwPolName string) error {
	dn := fmt.Sprintf("%s/rstoVpnGwPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnCloudVpnGwPolName": "%s","annotation": "orchestrator:terraform"
								
			}
		}
	}`, "cloudRsToVpnGwPol", dn, tnCloudVpnGwPolName))

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

func (sm *ServiceManager) ReadRelationcloudRsToVpnGwPolFromCloudVpnGateway(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "cloudRsToVpnGwPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "cloudRsToVpnGwPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationcloudRsToDirectConnPolFromCloudVpnGateway(parentDn, tnCloudDirectConnPolName string) error {
	dn := fmt.Sprintf("%s/rstoDirectConnPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnCloudDirectConnPolName": "%s","annotation": "orchestrator:terraform"
								
			}
		}
	}`, "cloudRsToDirectConnPol", dn, tnCloudDirectConnPolName))

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

func (sm *ServiceManager) ReadRelationcloudRsToDirectConnPolFromCloudVpnGateway(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "cloudRsToDirectConnPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "cloudRsToDirectConnPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationcloudRsToHostRouterPolFromCloudVpnGateway(parentDn, tnCloudHostRouterPolName string) error {
	dn := fmt.Sprintf("%s/rstoHostRouterPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnCloudHostRouterPolName": "%s","annotation": "orchestrator:terraform"
								
			}
		}
	}`, "cloudRsToHostRouterPol", dn, tnCloudHostRouterPolName))

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

func (sm *ServiceManager) ReadRelationcloudRsToHostRouterPolFromCloudVpnGateway(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "cloudRsToHostRouterPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "cloudRsToHostRouterPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
