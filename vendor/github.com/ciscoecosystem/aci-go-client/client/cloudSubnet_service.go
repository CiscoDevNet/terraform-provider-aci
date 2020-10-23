package client

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCloudSubnet(ip string, cloud_cidr_pool_dn string, description string, cloudSubnetattr models.CloudSubnetAttributes, zoneDn string) (*models.CloudSubnet, error) {
	rn := fmt.Sprintf("subnet-[%s]", ip)
	// parentDn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]", tenant, cloud_context_profile, cloud_cidr_pool_addr)
	parentDn := cloud_cidr_pool_dn
	cloudSubnet := models.NewCloudSubnet(rn, parentDn, description, cloudSubnetattr)
	jsonPayload, _, err := sm.PrepareModel(cloudSubnet)

	rsZoneAttachJSON := []byte(fmt.Sprintf(`
	{
		"cloudRsZoneAttach": {
			"attributes": {
				"annotation": "orchestrator:terraform",
				"dn": "%s/%s/rszoneAttach",
				"tDn": "%s"
			}
		}
	}
	`, parentDn, rn, zoneDn))
	zoneCon, err := container.ParseJSON(rsZoneAttachJSON)
	if err != nil {
		return nil, err
	}

	log.Printf("\n[DEBUG]asas %v", zoneCon.Data())
	jsonPayload.Array(cloudSubnet.ClassName, "children")
	jsonPayload.ArrayAppend(zoneCon.Data(), cloudSubnet.ClassName, "children")
	log.Printf("\n\n[DEBUG]asas %s\n\n", jsonPayload.String())
	jsonPayload.Set(ip, cloudSubnet.ClassName, "attributes", "ip")

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s/%s.json", parentDn, rn), jsonPayload, true)
	if err != nil {
		return nil, err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return nil, err
	}
	// err := sm.Save(cloudSubnet)
	return cloudSubnet, CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)
}

func (sm *ServiceManager) ReadCloudSubnet(ip string, cloud_cidr_pool_addr string, cloud_context_profile string, tenant string) (*models.CloudSubnet, error) {
	dn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]/subnet-[%s]", tenant, cloud_context_profile, cloud_cidr_pool_addr, ip)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudSubnet := models.CloudSubnetFromContainer(cont)
	return cloudSubnet, nil
}

func (sm *ServiceManager) DeleteCloudSubnet(ip string, cloud_cidr_pool_addr string, cloud_context_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]/subnet-[%s]", tenant, cloud_context_profile, cloud_cidr_pool_addr, ip)
	return sm.DeleteByDn(dn, models.CloudsubnetClassName)
}

func (sm *ServiceManager) UpdateCloudSubnet(ip string, cloud_cidr_pool_dn string, description string, cloudSubnetattr models.CloudSubnetAttributes, zoneDn string) (*models.CloudSubnet, error) {
	rn := fmt.Sprintf("subnet-[%s]", ip)
	// parentDn := fmt.Sprintf("uni/tn-%s/ctxprofile-%s/cidr-[%s]", tenant, cloud_context_profile, cloud_cidr_pool_addr)
	parentDn := cloud_cidr_pool_dn
	cloudSubnet := models.NewCloudSubnet(rn, parentDn, description, cloudSubnetattr)
	cloudSubnet.Status = "modified"

	jsonPayload, _, err := sm.PrepareModel(cloudSubnet)

	rsZoneAttachJSON := []byte(fmt.Sprintf(`
	{
		"cloudRsZoneAttach": {
			"attributes": {
				"annotation": "orchestrator:terraform",
				"dn": "%s/%s/rszoneAttach",
				"tDn": "%s"
			}
		}
	}
	`, parentDn, rn, zoneDn))
	zoneCon, err := container.ParseJSON(rsZoneAttachJSON)
	if err != nil {
		return nil, err
	}

	log.Printf("\n[DEBUG]asas %v", zoneCon.Data())
	jsonPayload.Array(cloudSubnet.ClassName, "children")
	jsonPayload.ArrayAppend(zoneCon.Data(), cloudSubnet.ClassName, "children")
	log.Printf("\n\n[DEBUG]asas %s\n\n", jsonPayload.String())
	jsonPayload.Set(ip, cloudSubnet.ClassName, "attributes", "ip")

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s/%s.json", parentDn, rn), jsonPayload, true)
	if err != nil {
		return nil, err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return nil, err
	}
	// err := sm.Save(cloudSubnet)
	return cloudSubnet, CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)

}

func (sm *ServiceManager) ListCloudSubnet(cloud_cidr_pool_addr string, cloud_context_profile string, tenant string) ([]*models.CloudSubnet, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctxprofile-%s/cidr-[%s]/cloudSubnet.json", baseurlStr, tenant, cloud_context_profile, cloud_cidr_pool_addr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudSubnetListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationcloudRsZoneAttachFromCloudSubnet(parentDn, tnCloudZoneName string) error {
	dn := fmt.Sprintf("%s/rszoneAttach", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "cloudRsZoneAttach", dn, tnCloudZoneName))

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

func (sm *ServiceManager) DeleteRelationcloudRsZoneAttachFromCloudSubnet(parentDn string) error {
	dn := fmt.Sprintf("%s/rszoneAttach", parentDn)
	return sm.DeleteByDn(dn, "cloudRsZoneAttach")
}

func (sm *ServiceManager) ReadRelationcloudRsZoneAttachFromCloudSubnet(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "cloudRsZoneAttach")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "cloudRsZoneAttach")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationcloudRsSubnetToFlowLogFromCloudSubnet(parentDn, tnCloudAwsFlowLogPolName string) error {
	dn := fmt.Sprintf("%s/rssubnetToFlowLog", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnCloudAwsFlowLogPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "cloudRsSubnetToFlowLog", dn, tnCloudAwsFlowLogPolName))

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

func (sm *ServiceManager) DeleteRelationcloudRsSubnetToFlowLogFromCloudSubnet(parentDn string) error {
	dn := fmt.Sprintf("%s/rssubnetToFlowLog", parentDn)
	return sm.DeleteByDn(dn, "cloudRsSubnetToFlowLog")
}

func (sm *ServiceManager) ReadRelationcloudRsSubnetToFlowLogFromCloudSubnet(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "cloudRsSubnetToFlowLog")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "cloudRsSubnetToFlowLog")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnCloudAwsFlowLogPolName")
		return dat, err
	} else {
		return nil, err
	}

}
