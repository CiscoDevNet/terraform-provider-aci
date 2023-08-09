package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreatePIMJPOutboundFilterPolicy(pim_interface_policy string, tenant string, description string, pimJPOutbFilterPolAttr models.PIMJPOutboundFilterPolicyAttributes) (*models.PIMJPOutboundFilterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimJPOutbFilterPol, tenant, pim_interface_policy)
	pimJPOutbFilterPol := models.NewPIMJPOutboundFilterPolicy(parentDn, description, pimJPOutbFilterPolAttr)

	err := sm.Save(pimJPOutbFilterPol)
	return pimJPOutbFilterPol, err
}

func (sm *ServiceManager) ReadPIMJPOutboundFilterPolicy(pim_interface_policy string, tenant string) (*models.PIMJPOutboundFilterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimJPOutbFilterPol, tenant, pim_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnPimJPOutbFilterPol)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	pimJPOutbFilterPol := models.PIMJPOutboundFilterPolicyFromContainer(cont)
	return pimJPOutbFilterPol, nil
}

func (sm *ServiceManager) DeletePIMJPOutboundFilterPolicy(pim_interface_policy string, tenant string) error {

	parentDn := fmt.Sprintf(models.ParentDnPimJPOutbFilterPol, tenant, pim_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnPimJPOutbFilterPol)

	return sm.DeleteByDn(dn, models.PimJPOutbFilterPolClassName)
}

func (sm *ServiceManager) UpdatePIMJPOutboundFilterPolicy(pim_interface_policy string, tenant string, description string, pimJPOutbFilterPolAttr models.PIMJPOutboundFilterPolicyAttributes) (*models.PIMJPOutboundFilterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimJPOutbFilterPol, tenant, pim_interface_policy)
	pimJPOutbFilterPol := models.NewPIMJPOutboundFilterPolicy(parentDn, description, pimJPOutbFilterPolAttr)

	pimJPOutbFilterPol.Status = "modified"
	err := sm.Save(pimJPOutbFilterPol)
	return pimJPOutbFilterPol, err
}

func (sm *ServiceManager) ListPIMJPOutboundFilterPolicy(pim_interface_policy string, tenant string) ([]*models.PIMJPOutboundFilterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimJPOutbFilterPol, tenant, pim_interface_policy)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.PimJPOutbFilterPolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.PIMJPOutboundFilterPolicyListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationPIMJPOutbFilterPolrtdmcRsFilterToRtMapPol(parentDn, tDn string) error {
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

func (sm *ServiceManager) DeleteRelationPIMJPOutbFilterPolrtdmcRsFilterToRtMapPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsfilterToRtMapPol", parentDn)
	return sm.DeleteByDn(dn, "rtdmcRsFilterToRtMapPol")
}

func (sm *ServiceManager) ReadRelationPIMJPOutbFilterPolrtdmcRsFilterToRtMapPol(parentDn string) (interface{}, error) {
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
