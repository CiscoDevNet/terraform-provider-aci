package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateLeafProfile(name string, description string, infraNodePattr models.LeafProfileAttributes) (*models.LeafProfile, error) {
	rn := fmt.Sprintf("infra/nprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraNodeP := models.NewLeafProfile(rn, parentDn, description, infraNodePattr)
	err := sm.Save(infraNodeP)
	return infraNodeP, err
}

func (sm *ServiceManager) ReadLeafProfile(name string) (*models.LeafProfile, error) {
	dn := fmt.Sprintf("uni/infra/nprof-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraNodeP := models.LeafProfileFromContainer(cont)
	return infraNodeP, nil
}

func (sm *ServiceManager) DeleteLeafProfile(name string) error {
	dn := fmt.Sprintf("uni/infra/nprof-%s", name)
	return sm.DeleteByDn(dn, models.InfranodepClassName)
}

func (sm *ServiceManager) UpdateLeafProfile(name string, description string, infraNodePattr models.LeafProfileAttributes) (*models.LeafProfile, error) {
	rn := fmt.Sprintf("infra/nprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraNodeP := models.NewLeafProfile(rn, parentDn, description, infraNodePattr)

	infraNodeP.Status = "modified"
	err := sm.Save(infraNodeP)
	return infraNodeP, err

}

func (sm *ServiceManager) ListLeafProfile() ([]*models.LeafProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infraNodeP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LeafProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsAccCardPFromLeafProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsaccCardP-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "infraRsAccCardP", tDn, dn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) DeleteRelationinfraRsAccCardPFromLeafProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsaccCardP-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "infraRsAccCardP")
}

func (sm *ServiceManager) ReadRelationinfraRsAccCardPFromLeafProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsAccCardP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsAccCardP")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationinfraRsAccPortPFromLeafProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsaccPortP-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "infraRsAccPortP", tDn, dn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) DeleteRelationinfraRsAccPortPFromLeafProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsaccPortP-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "infraRsAccPortP")
}

func (sm *ServiceManager) ReadRelationinfraRsAccPortPFromLeafProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsAccPortP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsAccPortP")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
