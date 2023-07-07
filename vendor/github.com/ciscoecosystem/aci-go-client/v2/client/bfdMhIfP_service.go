package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateAciBfdMultihopInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, nameAlias string, bfdMhIfPAttr models.AciBfdMultihopInterfaceProfileAttributes) (*models.AciBfdMultihopInterfaceProfile, error) {
	rn := fmt.Sprintf(models.RnbfdMhIfP)
	parentDn := fmt.Sprintf(models.ParentDnbfdMhIfP, tenant, l3_outside, logical_node_profile, logical_interface_profile)
	bfdMhIfP := models.NewAciBfdMultihopInterfaceProfile(rn, parentDn, description, nameAlias, bfdMhIfPAttr)
	err := sm.Save(bfdMhIfP)
	return bfdMhIfP, err
}

func (sm *ServiceManager) ReadAciBfdMultihopInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.AciBfdMultihopInterfaceProfile, error) {
	dn := fmt.Sprintf(models.DnbfdMhIfP, tenant, l3_outside, logical_node_profile, logical_interface_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bfdMhIfP := models.AciBfdMultihopInterfaceProfileFromContainer(cont)
	return bfdMhIfP, nil
}

func (sm *ServiceManager) DeleteAciBfdMultihopInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf(models.DnbfdMhIfP, tenant, l3_outside, logical_node_profile, logical_interface_profile)
	return sm.DeleteByDn(dn, models.BfdmhifpClassName)
}

func (sm *ServiceManager) UpdateAciBfdMultihopInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, nameAlias string, bfdMhIfPAttr models.AciBfdMultihopInterfaceProfileAttributes) (*models.AciBfdMultihopInterfaceProfile, error) {
	rn := fmt.Sprintf(models.RnbfdMhIfP)
	parentDn := fmt.Sprintf(models.ParentDnbfdMhIfP, tenant, l3_outside, logical_node_profile, logical_interface_profile)
	bfdMhIfP := models.NewAciBfdMultihopInterfaceProfile(rn, parentDn, description, nameAlias, bfdMhIfPAttr)
	bfdMhIfP.Status = "modified"
	err := sm.Save(bfdMhIfP)
	return bfdMhIfP, err
}

func (sm *ServiceManager) ListAciBfdMultihopInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.AciBfdMultihopInterfaceProfile, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/bfdMhIfP.json", models.BaseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.AciBfdMultihopInterfaceProfileListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationbfdRsMhIfPol(parentDn, annotation, tnBfdMhIfPolName string) error {
	dn := fmt.Sprintf("%s/rsMhIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnBfdMhIfPolName": "%s"
			}
		}
	}`, "bfdRsMhIfPol", dn, annotation, tnBfdMhIfPolName))

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

func (sm *ServiceManager) DeleteRelationbfdRsMhIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsMhIfPol", parentDn)
	return sm.DeleteByDn(dn, "bfdRsMhIfPol")
}

func (sm *ServiceManager) ReadRelationbfdRsMhIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "bfdRsMhIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "bfdRsMhIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnBfdMhIfPolName")
		return dat, err
	} else {
		return nil, err
	}
}
