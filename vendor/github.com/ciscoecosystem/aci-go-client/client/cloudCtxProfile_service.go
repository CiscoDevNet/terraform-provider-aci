package client

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudContextProfile(name string, tenant string, description string, cloudCtxProfileattr models.CloudContextProfileAttributes, primaryCidr string, region, vendor string, vrf string) (*models.CloudContextProfile, error) {
	rn := fmt.Sprintf("ctxprofile-%s", name)
	parentDn := tenant
	cloudCtxProfile := models.NewCloudContextProfile(rn, parentDn, description, cloudCtxProfileattr)
	jsonPayload, _, err := sm.PrepareModel(cloudCtxProfile)

	cidrJSON := []byte(fmt.Sprintf(`
	{
		"cloudCidr": {
			"attributes": {
				"addr": "%s",
				"primary": "yes",
				"status": "created,modified"
			}
		}
	}
	`, primaryCidr))

	regionAttach := []byte(fmt.Sprintf(`
	{
		"cloudRsCtxProfileToRegion": {
			"attributes": {
				"status": "created,modified",
				"tDn": "uni/clouddomp/provp-%s/region-%s"
			}
		}
	}
	`, vendor, region))

	ctxAttach := []byte(fmt.Sprintf(`
	{
		"cloudRsToCtx": {
			"attributes": {
				"tnFvCtxName": "%s"
			}
		}
	}
	`, vrf))

	cidrCon, err := container.ParseJSON(cidrJSON)
	regionCon, err := container.ParseJSON(regionAttach)
	vrfCon, err := container.ParseJSON(ctxAttach)
	if err != nil {
		return nil, err
	}
	if err != nil {

	}

	log.Printf("\n\n\n[DEBUG]nknk %v", vrfCon.Data())
	jsonPayload.Array(cloudCtxProfile.ClassName, "children")
	jsonPayload.ArrayAppend(vrfCon.Data(), cloudCtxProfile.ClassName, "children")

	jsonPayload.ArrayAppend(cidrCon.Data(), cloudCtxProfile.ClassName, "children")
	jsonPayload.ArrayAppend(regionCon.Data(), cloudCtxProfile.ClassName, "children")

	log.Printf("\n\n\n\n[DEBUG]nkdemo%s\n\n\n\n", jsonPayload.String())
	jsonPayload.Set(name, cloudCtxProfile.ClassName, "attributes", "name")
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s/%s.json", parentDn, rn), jsonPayload, true)
	if err != nil {
		return nil, err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return nil, err
	}

	//err := sm.Save(cloudCtxProfile)
	return cloudCtxProfile, CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)
}

func (sm *ServiceManager) ReadCloudContextProfile(name string, tenant string) (*models.CloudContextProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudCtxProfile := models.CloudContextProfileFromContainer(cont)
	return cloudCtxProfile, nil
}

func (sm *ServiceManager) DeleteCloudContextProfile(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s", tenant, name)
	return sm.DeleteByDn(dn, models.CloudctxprofileClassName)
}

func (sm *ServiceManager) UpdateCloudContextProfile(name string, tenant string, description string, cloudCtxProfileattr models.CloudContextProfileAttributes, primaryCidr string, region, vendor string, vrf string) (*models.CloudContextProfile, error) {
	rn := fmt.Sprintf("ctxprofile-%s", name)
	parentDn := tenant
	cloudCtxProfile := models.NewCloudContextProfile(rn, parentDn, description, cloudCtxProfileattr)
	cloudCtxProfile.Status = "modified"
	jsonPayload, _, err := sm.PrepareModel(cloudCtxProfile)

	cidrJSON := []byte(fmt.Sprintf(`
	{
		"cloudCidr": {
			"attributes": {
				"addr": "%s",
				"primary": "yes",
				"status": "modified"
			}
		}
	}
	`, primaryCidr))

	regionAttach := []byte(fmt.Sprintf(`
	{
		"cloudRsCtxProfileToRegion": {
			"attributes": {
				"status": "modified",
				"tDn": "uni/clouddomp/provp-%s/region-%s"
			}
		}
	}
	`, vendor, region))

	ctxAttach := []byte(fmt.Sprintf(`
	{
		"cloudRsToCtx": {
			"attributes": {
				"tnFvCtxName": "%s",
			}
		}
	}
	`, vrf))

	cidrCon, err := container.ParseJSON(cidrJSON)
	regionCon, err := container.ParseJSON(regionAttach)
	vrfCon, err := container.ParseJSON(ctxAttach)
	jsonPayload.Array(cloudCtxProfile.ClassName, "children")
	jsonPayload.ArrayAppend(cidrCon.Data(), cloudCtxProfile.ClassName, "children")
	jsonPayload.ArrayAppend(regionCon.Data(), cloudCtxProfile.ClassName, "children")
	jsonPayload.ArrayAppend(vrfCon.Data(), cloudCtxProfile.ClassName, "children")

	jsonPayload.Set(name, cloudCtxProfile.ClassName, "attributes", "name")
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s/%s.json", parentDn, rn), jsonPayload, true)
	if err != nil {
		return nil, err
	}

	cont, _, err := sm.client.Do(req)
	if err != nil {
		return nil, err
	}

	return cloudCtxProfile, CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)

}

func (sm *ServiceManager) ListCloudContextProfile(tenant string) ([]*models.CloudContextProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/cloudCtxProfile.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudContextProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationcloudRsCtxToFlowLogFromCloudContextProfile(parentDn, tnCloudAwsFlowLogPolName string) error {
	dn := fmt.Sprintf("%s/rsctxToFlowLog", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnCloudAwsFlowLogPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "cloudRsCtxToFlowLog", dn, tnCloudAwsFlowLogPolName))

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

func (sm *ServiceManager) DeleteRelationcloudRsCtxToFlowLogFromCloudContextProfile(parentDn string) error {
	dn := fmt.Sprintf("%s/rsctxToFlowLog", parentDn)
	return sm.DeleteByDn(dn, "cloudRsCtxToFlowLog")
}
func (sm *ServiceManager) CreateRelationcloudRsToCtxFromCloudContextProfile(parentDn, tnFvCtxName string) error {
	dn := fmt.Sprintf("%s/rstoCtx", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvCtxName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "cloudRsToCtx", dn, tnFvCtxName))

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
func (sm *ServiceManager) CreateRelationcloudRsCtxProfileToRegionFromCloudContextProfile(parentDn, tnCloudRegionName string) error {
	dn := fmt.Sprintf("%s/rsctxProfileToRegion", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "cloudRsCtxProfileToRegion", dn, tnCloudRegionName))

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

func (sm *ServiceManager) DeleteRelationcloudRsCtxProfileToRegionFromCloudContextProfile(parentDn string) error {
	dn := fmt.Sprintf("%s/rsctxProfileToRegion", parentDn)
	return sm.DeleteByDn(dn, "cloudRsCtxProfileToRegion")
}

func (sm *ServiceManager) CreateRelationcloudRsCtxProfileTocloudRsCtxProfileToGatewayRouterP(parentDn, tDN string) error {
	dn := fmt.Sprintf("%s/rsctxProfileToGatewayRouterP-[%s]", parentDn, tDN)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "cloudRsCtxProfileToGatewayRouterP", dn, tDN))

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

func (sm *ServiceManager) DeleteRelationcloudRsCtxProfileTocloudRsCtxProfileToGatewayRouterP(parentDn, tDN string) error {
	dn := fmt.Sprintf("%s/rsctxProfileToGatewayRouterP-[%s]", parentDn, tDN)
	return sm.DeleteByDn(dn, "cloudRsCtxProfileToGatewayRouterP")
}
