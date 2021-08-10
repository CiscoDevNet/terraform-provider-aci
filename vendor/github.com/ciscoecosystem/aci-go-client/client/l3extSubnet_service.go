package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3ExtSubnet(ip string, external_network_instance_profile string, l3_outside string, tenant string, description string, l3extSubnetattr models.L3ExtSubnetAttributes) (*models.L3ExtSubnet, error) {
	rn := fmt.Sprintf("extsubnet-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/instP-%s", tenant, l3_outside, external_network_instance_profile)
	l3extSubnet := models.NewL3ExtSubnet(rn, parentDn, description, l3extSubnetattr)
	err := sm.Save(l3extSubnet)
	return l3extSubnet, err
}

func (sm *ServiceManager) ReadL3ExtSubnet(ip string, external_network_instance_profile string, l3_outside string, tenant string) (*models.L3ExtSubnet, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/instP-%s/extsubnet-[%s]", tenant, l3_outside, external_network_instance_profile, ip)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extSubnet := models.L3ExtSubnetFromContainer(cont)
	return l3extSubnet, nil
}

func (sm *ServiceManager) DeleteL3ExtSubnet(ip string, external_network_instance_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/instP-%s/extsubnet-[%s]", tenant, l3_outside, external_network_instance_profile, ip)
	return sm.DeleteByDn(dn, models.L3extsubnetClassName)
}

func (sm *ServiceManager) UpdateL3ExtSubnet(ip string, external_network_instance_profile string, l3_outside string, tenant string, description string, l3extSubnetattr models.L3ExtSubnetAttributes) (*models.L3ExtSubnet, error) {
	rn := fmt.Sprintf("extsubnet-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/instP-%s", tenant, l3_outside, external_network_instance_profile)
	l3extSubnet := models.NewL3ExtSubnet(rn, parentDn, description, l3extSubnetattr)

	l3extSubnet.Status = "modified"
	err := sm.Save(l3extSubnet)
	return l3extSubnet, err

}

func (sm *ServiceManager) ListL3L3ExtSubnet(external_network_instance_profile string, l3_outside string, tenant string) ([]*models.L3ExtSubnet, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/instP-%s/l3extSubnet.json", baseurlStr, tenant, l3_outside, external_network_instance_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3ExtSubnetListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationl3extRsSubnetToProfileFromL3ExtSubnet(parentDn, tnRtctrlProfileName, direction string) error {
	dn := fmt.Sprintf("%s/rssubnetToProfile-[%s]-%s", parentDn, tnRtctrlProfileName, direction)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "l3extRsSubnetToProfile", dn))

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

func (sm *ServiceManager) DeleteRelationl3extRsSubnetToProfileFromL3ExtSubnet(parentDn, tnRtctrlProfileName, direction string) error {
	dn := fmt.Sprintf("%s/rssubnetToProfile-[%s]-%s", parentDn, tnRtctrlProfileName, direction)
	return sm.DeleteByDn(dn, "l3extRsSubnetToProfile")
}

func (sm *ServiceManager) ReadRelationl3extRsSubnetToProfileFromL3ExtSubnet(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsSubnetToProfile")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsSubnetToProfile")

	st := make([]map[string]string, 0)

	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnRtctrlProfileName"] = models.G(contItem, "tDn")
		paramMap["direction"] = models.G(contItem, "direction")

		st = append(st, paramMap)

	}

	return st, err

}
func (sm *ServiceManager) CreateRelationl3extRsSubnetToRtSummFromL3ExtSubnet(parentDn, tnRtsumARtSummPolName string) error {
	dn := fmt.Sprintf("%s/rsSubnetToRtSumm", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsSubnetToRtSumm", dn, tnRtsumARtSummPolName))

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

func (sm *ServiceManager) DeleteRelationl3extRsSubnetToRtSummFromL3ExtSubnet(parentDn string) error {
	dn := fmt.Sprintf("%s/rsSubnetToRtSumm", parentDn)
	return sm.DeleteByDn(dn, "l3extRsSubnetToRtSumm")
}

func (sm *ServiceManager) ReadRelationl3extRsSubnetToRtSummFromL3ExtSubnet(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsSubnetToRtSumm")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsSubnetToRtSumm")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
