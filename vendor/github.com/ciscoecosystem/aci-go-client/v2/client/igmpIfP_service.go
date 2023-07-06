package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateIGMPInterfaceProfile(parentDn string, description string, igmpIfPAttr models.IGMPInterfaceProfileAttributes) (*models.IGMPInterfaceProfile, error) {

	igmpIfP := models.NewIGMPInterfaceProfile(parentDn, description, igmpIfPAttr)

	err := sm.Save(igmpIfP)
	return igmpIfP, err
}

func (sm *ServiceManager) ReadIGMPInterfaceProfile(parentDn string) (*models.InterfaceProfile, error) {

	dn := fmt.Sprintf("%s/%s", parentDn, models.RnIgmpIfP)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	igmpIfP := models.InterfaceProfileFromContainer(cont)
	return igmpIfP, nil
}

func (sm *ServiceManager) DeleteIGMPInterfaceProfile(parentDn string) error {

	dn := fmt.Sprintf("%s/%s", parentDn, models.RnIgmpIfP)

	return sm.DeleteByDn(dn, models.IgmpIfPClassName)
}

func (sm *ServiceManager) UpdateIGMPInterfaceProfile(parentDn string, description string, igmpIfPAttr models.IGMPInterfaceProfileAttributes) (*models.IGMPInterfaceProfile, error) {

	igmpIfP := models.NewIGMPInterfaceProfile(parentDn, description, igmpIfPAttr)

	igmpIfP.Status = "modified"
	err := sm.Save(igmpIfP)
	return igmpIfP, err
}

func (sm *ServiceManager) ListIGMPInterfaceProfile(parentDn string) ([]*models.InterfaceProfile, error) {

	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.IgmpIfPClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.InterfaceProfileListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationIGMPRsIfPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"	
			}
		}
	}`, "igmpRsIfPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationIGMPRsIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsIfPol", parentDn)
	return sm.DeleteByDn(dn, "igmpRsIfPol")
}

func (sm *ServiceManager) ReadRelationIGMPRsIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "igmpRsIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "igmpRsIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
