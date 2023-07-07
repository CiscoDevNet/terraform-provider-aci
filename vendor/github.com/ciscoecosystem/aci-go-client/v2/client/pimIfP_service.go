package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreatePIMInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, pimIfPAttr models.PIMInterfaceProfileAttributes) (*models.PIMInterfaceProfile, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimIfP, tenant, l3_outside, logical_node_profile, logical_interface_profile)
	pimIfP := models.NewPIMInterfaceProfile(parentDn, description, pimIfPAttr)

	err := sm.Save(pimIfP)
	return pimIfP, err
}

func (sm *ServiceManager) ReadPIMInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.PIMInterfaceProfile, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimIfP, tenant, l3_outside, logical_node_profile, logical_interface_profile)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnPimIfP)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	pimIfP := models.PIMInterfaceProfileFromContainer(cont)
	return pimIfP, nil
}

func (sm *ServiceManager) DeletePIMInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {

	parentDn := fmt.Sprintf(models.ParentDnPimIfP, tenant, l3_outside, logical_node_profile, logical_interface_profile)
	dn := fmt.Sprintf("%s/%s", parentDn, models.RnPimIfP)

	return sm.DeleteByDn(dn, models.PimIfPClassName)
}

func (sm *ServiceManager) UpdatePIMInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, pimIfPAttr models.PIMInterfaceProfileAttributes) (*models.PIMInterfaceProfile, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimIfP, tenant, l3_outside, logical_node_profile, logical_interface_profile)
	pimIfP := models.NewPIMInterfaceProfile(parentDn, description, pimIfPAttr)

	pimIfP.Status = "modified"
	err := sm.Save(pimIfP)
	return pimIfP, err
}

func (sm *ServiceManager) ListPIMInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.PIMInterfaceProfile, error) {

	parentDn := fmt.Sprintf(models.ParentDnPimIfP, tenant, l3_outside, logical_node_profile, logical_interface_profile)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.PimIfPClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.PIMInterfaceProfileListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationPIMRsIfPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"	
			}
		}
	}`, "pimRsIfPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationPIMRsIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsIfPol", parentDn)
	return sm.DeleteByDn(dn, "pimRsIfPol")
}

func (sm *ServiceManager) ReadRelationPIMRsIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "pimRsIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "pimRsIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
