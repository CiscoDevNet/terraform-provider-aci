package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateIGMPReportPolicy(igmp_interface_policy string, tenant string, description string, igmpRepPolAttr models.IGMPReportPolicyAttributes) (*models.IGMPReportPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpRepPol, tenant, igmp_interface_policy)
	igmpRepPol := models.NewIGMPReportPolicy(parentDn, description, igmpRepPolAttr)

	err := sm.Save(igmpRepPol)
	return igmpRepPol, err
}

func (sm *ServiceManager) ReadIGMPReportPolicy(igmp_interface_policy string, tenant string) (*models.IGMPReportPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpRepPol, tenant, igmp_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnIgmpRepPol)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	igmpRepPol := models.IGMPReportPolicyFromContainer(cont)
	return igmpRepPol, nil
}

func (sm *ServiceManager) DeleteIGMPReportPolicy(igmp_interface_policy string, tenant string) error {

	parentDn := fmt.Sprintf(models.ParentDnIgmpRepPol, tenant, igmp_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnIgmpRepPol)

	return sm.DeleteByDn(dn, models.IgmpRepPolClassName)
}

func (sm *ServiceManager) UpdateIGMPReportPolicy(igmp_interface_policy string, tenant string, description string, igmpRepPolAttr models.IGMPReportPolicyAttributes) (*models.IGMPReportPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpRepPol, tenant, igmp_interface_policy)
	igmpRepPol := models.NewIGMPReportPolicy(parentDn, description, igmpRepPolAttr)

	igmpRepPol.Status = "modified"
	err := sm.Save(igmpRepPol)
	return igmpRepPol, err
}

func (sm *ServiceManager) ListIGMPReportPolicy(igmp_interface_policy string, tenant string) ([]*models.IGMPReportPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnIgmpRepPol, tenant, igmp_interface_policy)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.IgmpRepPolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.IGMPReportPolicyListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationigmpRepPolrtdmcRsFilterToRtMapPol(parentDn, tDn string) error {
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

func (sm *ServiceManager) DeleteRelationigmpRepPolrtdmcRsFilterToRtMapPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsfilterToRtMapPol", parentDn)
	return sm.DeleteByDn(dn, "rtdmcRsFilterToRtMapPol")
}

func (sm *ServiceManager) ReadRelationigmpRepPolrtdmcRsFilterToRtMapPol(parentDn string) (interface{}, error) {
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
