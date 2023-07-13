package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateIGMPStateLimitPolicy(igmp_interface_policy string, tenant string, description string, igmpStateLPolAttr models.IGMPStateLimitPolicyAttributes) (*models.IGMPStateLimitPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpStateLPol, tenant, igmp_interface_policy)
	igmpStateLPol := models.NewIGMPStateLimitPolicy(parentDn, description, igmpStateLPolAttr)

	err := sm.Save(igmpStateLPol)
	return igmpStateLPol, err
}

func (sm *ServiceManager) ReadIGMPStateLimitPolicy(igmp_interface_policy string, tenant string) (*models.IGMPStateLimitPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpStateLPol, tenant, igmp_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnIgmpStateLPol)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	igmpStateLPol := models.IGMPStateLimitPolicyFromContainer(cont)
	return igmpStateLPol, nil
}

func (sm *ServiceManager) DeleteIGMPStateLimitPolicy(igmp_interface_policy string, tenant string) error {

	parentDn := fmt.Sprintf(models.ParentDnIgmpStateLPol, tenant, igmp_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnIgmpStateLPol)

	return sm.DeleteByDn(dn, models.IgmpStateLPolClassName)
}

func (sm *ServiceManager) UpdateIGMPStateLimitPolicy(igmp_interface_policy string, tenant string, description string, igmpStateLPolAttr models.IGMPStateLimitPolicyAttributes) (*models.IGMPStateLimitPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpStateLPol, tenant, igmp_interface_policy)
	igmpStateLPol := models.NewIGMPStateLimitPolicy(parentDn, description, igmpStateLPolAttr)

	igmpStateLPol.Status = "modified"
	err := sm.Save(igmpStateLPol)
	return igmpStateLPol, err
}

func (sm *ServiceManager) ListIGMPStateLimitPolicy(igmp_interface_policy string, tenant string) ([]*models.IGMPStateLimitPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpStateLPol, tenant, igmp_interface_policy)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.IgmpStateLPolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.IGMPStateLimitPolicyListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationigmpStateLPolrtdmcRsFilterToRtMapPol(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsfilterToRtMapPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "orchestrator:terraform",
				"tDn": "%s"	
			}
		}
	}`, "rtdmcRsFilterToRtMapPol", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationigmpStateLPolrtdmcRsFilterToRtMapPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsfilterToRtMapPol", parentDn)
	return sm.DeleteByDn(dn, "rtdmcRsFilterToRtMapPol")
}

func (sm *ServiceManager) ReadRelationigmpStateLPolrtdmcRsFilterToRtMapPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "rtdmcRsFilterToRtMapPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "rtdmcRsFilterToRtMapPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
