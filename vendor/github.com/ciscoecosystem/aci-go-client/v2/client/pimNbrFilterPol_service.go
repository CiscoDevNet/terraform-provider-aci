package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreatePIMNeighborFiterPolicy(pim_interface_policy string, tenant string, description string, pimNbrFilterPolAttr models.PIMNeighborFiterPolicyAttributes) (*models.PIMNeighborFiterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimNbrFilterPol, tenant, pim_interface_policy)
	pimNbrFilterPol := models.NewPIMNeighborFiterPolicy(parentDn, description, pimNbrFilterPolAttr)

	err := sm.Save(pimNbrFilterPol)
	return pimNbrFilterPol, err
}

func (sm *ServiceManager) ReadPIMNeighborFiterPolicy(pim_interface_policy string, tenant string) (*models.PIMNeighborFiterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimNbrFilterPol, tenant, pim_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnPimNbrFilterPol)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	pimNbrFilterPol := models.PIMNeighborFiterPolicyFromContainer(cont)
	return pimNbrFilterPol, nil
}

func (sm *ServiceManager) DeletePIMNeighborFiterPolicy(pim_interface_policy string, tenant string) error {

	parentDn := fmt.Sprintf(models.ParentDnPimNbrFilterPol, tenant, pim_interface_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnPimNbrFilterPol)

	return sm.DeleteByDn(dn, models.PimNbrFilterPolClassName)
}

func (sm *ServiceManager) UpdatePIMNeighborFiterPolicy(pim_interface_policy string, tenant string, description string, pimNbrFilterPolAttr models.PIMNeighborFiterPolicyAttributes) (*models.PIMNeighborFiterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimNbrFilterPol, tenant, pim_interface_policy)
	pimNbrFilterPol := models.NewPIMNeighborFiterPolicy(parentDn, description, pimNbrFilterPolAttr)

	pimNbrFilterPol.Status = "modified"
	err := sm.Save(pimNbrFilterPol)
	return pimNbrFilterPol, err
}

func (sm *ServiceManager) ListPIMNeighborFiterPolicy(pim_interface_policy string, tenant string) ([]*models.PIMNeighborFiterPolicy, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimNbrFilterPol, tenant, pim_interface_policy)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.PimNbrFilterPolClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.PIMNeighborFiterPolicyListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationPIMNbrFilterPolrtdmcRsFilterToRtMapPol(parentDn, tDn string) error {
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

func (sm *ServiceManager) DeleteRelationPIMNbrFilterPolrtdmcRsFilterToRtMapPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsfilterToRtMapPol", parentDn)
	return sm.DeleteByDn(dn, "rtdmcRsFilterToRtMapPol")
}

func (sm *ServiceManager) ReadRelationPIMNbrFilterPolrtdmcRsFilterToRtMapPol(parentDn string) (interface{}, error) {
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
