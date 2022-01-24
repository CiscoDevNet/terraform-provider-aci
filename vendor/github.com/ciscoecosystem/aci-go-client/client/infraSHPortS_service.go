package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSpineAccessPortSelector(spine_access_port_selector_type string, name string, spine_interface_profile string, description string, nameAlias string, infraSHPortSAttr models.SpineAccessPortSelectorAttributes) (*models.SpineAccessPortSelector, error) {
	rn := fmt.Sprintf(models.RninfraSHPortS, name, spine_access_port_selector_type)
	parentDn := fmt.Sprintf(models.ParentDninfraSHPortS, spine_interface_profile)
	infraSHPortS := models.NewSpineAccessPortSelector(rn, parentDn, description, nameAlias, infraSHPortSAttr)
	err := sm.Save(infraSHPortS)
	return infraSHPortS, err
}

func (sm *ServiceManager) ReadSpineAccessPortSelector(spine_access_port_selector_type string, name string, spine_interface_profile string) (*models.SpineAccessPortSelector, error) {
	dn := fmt.Sprintf(models.DninfraSHPortS, spine_interface_profile, name, spine_access_port_selector_type)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	infraSHPortS := models.SpineAccessPortSelectorFromContainer(cont)
	return infraSHPortS, nil
}

func (sm *ServiceManager) DeleteSpineAccessPortSelector(spine_access_port_selector_type string, name string, spine_interface_profile string) error {
	dn := fmt.Sprintf(models.DninfraSHPortS, spine_interface_profile, name, spine_access_port_selector_type)
	return sm.DeleteByDn(dn, models.InfrashportsClassName)
}

func (sm *ServiceManager) UpdateSpineAccessPortSelector(spine_access_port_selector_type string, name string, spine_interface_profile string, description string, nameAlias string, infraSHPortSAttr models.SpineAccessPortSelectorAttributes) (*models.SpineAccessPortSelector, error) {
	rn := fmt.Sprintf(models.RninfraSHPortS, name, spine_access_port_selector_type)
	parentDn := fmt.Sprintf(models.ParentDninfraSHPortS, spine_interface_profile)
	infraSHPortS := models.NewSpineAccessPortSelector(rn, parentDn, description, nameAlias, infraSHPortSAttr)
	infraSHPortS.Status = "modified"
	err := sm.Save(infraSHPortS)
	return infraSHPortS, err
}

func (sm *ServiceManager) ListSpineAccessPortSelector(spine_interface_profile string) ([]*models.SpineAccessPortSelector, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/spaccportprof-%s/infraSHPortS.json", models.BaseurlStr, spine_interface_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SpineAccessPortSelectorListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsSpAccGrp(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsspAccGrp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "infraRsSpAccGrp", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationinfraRsSpAccGrp(parentDn string) error {
	dn := fmt.Sprintf("%s/rsspAccGrp", parentDn)
	return sm.DeleteByDn(dn, "infraRsSpAccGrp")
}

func (sm *ServiceManager) ReadRelationinfraRsSpAccGrp(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsSpAccGrp")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsSpAccGrp")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
