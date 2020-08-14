package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func (sm *ServiceManager) CreateSpineProfile(name string, description string, infraSpinePattr models.SpineProfileAttributes) (*models.SpineProfile, error) {
	rn := fmt.Sprintf("infra/spprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraSpineP := models.NewSpineProfile(rn, parentDn, description, infraSpinePattr)
	err := sm.Save(infraSpineP)
	return infraSpineP, err
}

func (sm *ServiceManager) ReadSpineProfile(name string) (*models.SpineProfile, error) {
	dn := fmt.Sprintf("uni/infra/spprof-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpineP := models.SpineProfileFromContainer(cont)
	return infraSpineP, nil
}

func (sm *ServiceManager) DeleteSpineProfile(name string) error {
	dn := fmt.Sprintf("uni/infra/spprof-%s", name)
	return sm.DeleteByDn(dn, models.InfraspinepClassName)
}

func (sm *ServiceManager) UpdateSpineProfile(name string, description string, infraSpinePattr models.SpineProfileAttributes) (*models.SpineProfile, error) {
	rn := fmt.Sprintf("infra/spprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraSpineP := models.NewSpineProfile(rn, parentDn, description, infraSpinePattr)

	infraSpineP.Status = "modified"
	err := sm.Save(infraSpineP)
	return infraSpineP, err

}

func (sm *ServiceManager) ListSpineProfile() ([]*models.SpineProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infraSpineP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.SpineProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsSpAccPortPFromSpineProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsspAccPortP-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "infraRsSpAccPortP", dn))

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

func (sm *ServiceManager) DeleteRelationinfraRsSpAccPortPFromSpineProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsspAccPortP-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "infraRsSpAccPortP")
}

func (sm *ServiceManager) ReadRelationinfraRsSpAccPortPFromSpineProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsSpAccPortP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsSpAccPortP")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
