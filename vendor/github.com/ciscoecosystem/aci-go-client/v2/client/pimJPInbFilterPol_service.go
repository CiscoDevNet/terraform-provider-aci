package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreatePIMJPInboundFilterPolicy(pim_interface_policy string, tenant string, description string, pimJPInbFilterPolAttr models.PIMJPInboundFilterPolicyAttributes) (*models.PIMJPInboundFilterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimJPInbFilterPol, tenant, pim_interface_policy)
	pimJPInbFilterPol := models.NewPIMJPInboundFilterPolicy(parentDn, description, pimJPInbFilterPolAttr)

	err := sm.Save(pimJPInbFilterPol)
	return pimJPInbFilterPol, err
}

func (sm *ServiceManager) ReadPIMJPInboundFilterPolicy(pim_interface_policy string, tenant string) (*models.PIMJPInboundFilterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimJPInbFilterPol, tenant, pim_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnPimJPInbFilterPol)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	pimJPInbFilterPol := models.PIMJPInboundFilterPolicyFromContainer(cont)
	return pimJPInbFilterPol, nil
}

func (sm *ServiceManager) DeletePIMJPInboundFilterPolicy(pim_interface_policy string, tenant string) error {

	parentDn := fmt.Sprintf(models.ParentDnPimJPInbFilterPol, tenant, pim_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnPimJPInbFilterPol)

	return sm.DeleteByDn(dn, models.PimJPInbFilterPolClassName)
}

func (sm *ServiceManager) UpdatePIMJPInboundFilterPolicy(pim_interface_policy string, tenant string, description string, pimJPInbFilterPolAttr models.PIMJPInboundFilterPolicyAttributes) (*models.PIMJPInboundFilterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimJPInbFilterPol, tenant, pim_interface_policy)
	pimJPInbFilterPol := models.NewPIMJPInboundFilterPolicy(parentDn, description, pimJPInbFilterPolAttr)

	pimJPInbFilterPol.Status = "modified"
	err := sm.Save(pimJPInbFilterPol)
	return pimJPInbFilterPol, err
}

func (sm *ServiceManager) ListPIMJPInboundFilterPolicy(pim_interface_policy string, tenant string) ([]*models.PIMJPInboundFilterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimJPInbFilterPol, tenant, pim_interface_policy)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.PimJPInbFilterPolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.PIMJPInboundFilterPolicyListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationPIMJPInbFilterPolrtdmcRsFilterToRtMapPol(parentDn, tDn string) error {
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

func (sm *ServiceManager) DeleteRelationPIMJPInbFilterPolrtdmcRsFilterToRtMapPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsfilterToRtMapPol", parentDn)
	return sm.DeleteByDn(dn, "rtdmcRsFilterToRtMapPol")
}

func (sm *ServiceManager) ReadRelationPIMJPInbFilterPolrtdmcRsFilterToRtMapPol(parentDn string) (interface{}, error) {
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
