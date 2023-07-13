package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateIGMPStaticReportPolicy(igmp_interface_policy string, tenant string, description string, igmpStRepPolAttr models.IGMPStaticReportPolicyAttributes) (*models.IGMPStaticReportPolicy, error) {
	parentDn := fmt.Sprintf(models.ParentDnIgmpStRepPol, tenant, igmp_interface_policy)
	igmpStRepPol := models.NewIGMPStaticReportPolicy(parentDn, description, igmpStRepPolAttr)

	err := sm.Save(igmpStRepPol)
	return igmpStRepPol, err
}

func (sm *ServiceManager) ReadIGMPStaticReportPolicy(igmp_interface_policy string, tenant string) (*models.IGMPStaticReportPolicy, error) {
	parentDn := fmt.Sprintf(models.ParentDnIgmpStRepPol, tenant, igmp_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnIgmpStRepPol)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	igmpStRepPol := models.IGMPStaticReportPolicyFromContainer(cont)
	return igmpStRepPol, nil
}

func (sm *ServiceManager) DeleteIGMPStaticReportPolicy(igmp_interface_policy string, tenant string) error {

	parentDn := fmt.Sprintf(models.ParentDnIgmpStRepPol, tenant, igmp_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnIgmpStRepPol)

	return sm.DeleteByDn(dn, models.IgmpStRepPolClassName)
}

func (sm *ServiceManager) UpdateIGMPStaticReportPolicy(igmp_interface_policy string, tenant string, description string, igmpStRepPolAttr models.IGMPStaticReportPolicyAttributes) (*models.IGMPStaticReportPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpStRepPol, tenant, igmp_interface_policy)
	igmpStRepPol := models.NewIGMPStaticReportPolicy(parentDn, description, igmpStRepPolAttr)

	igmpStRepPol.Status = "modified"
	err := sm.Save(igmpStRepPol)
	return igmpStRepPol, err
}

func (sm *ServiceManager) ListIGMPStaticReportPolicy(igmp_interface_policy string, tenant string) ([]*models.IGMPStaticReportPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpStRepPol, tenant, igmp_interface_policy)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.IgmpStRepPolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.IGMPStaticReportPolicyListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationigmpStRepPolrtdmcRsFilterToRtMapPol(parentDn, tDn string) error {
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

func (sm *ServiceManager) DeleteRelationigmpStRepPolrtdmcRsFilterToRtMapPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsfilterToRtMapPol", parentDn)
	return sm.DeleteByDn(dn, "rtdmcRsFilterToRtMapPol")
}

func (sm *ServiceManager) ReadRelationigmpStRepPolrtdmcRsFilterToRtMapPol(parentDn string) (interface{}, error) {
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
