package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateAny(vrf string, tenant string, description string, vzAnyattr models.AnyAttributes) (*models.Any, error) {
	rn := fmt.Sprintf("any")
	parentDn := fmt.Sprintf("uni/tn-%s/ctx-%s", tenant, vrf)
	vzAny := models.NewAny(rn, parentDn, description, vzAnyattr)
	err := sm.Save(vzAny)
	return vzAny, err
}

func (sm *ServiceManager) ReadAny(vrf string, tenant string) (*models.Any, error) {
	dn := fmt.Sprintf("uni/tn-%s/ctx-%s/any", tenant, vrf)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzAny := models.AnyFromContainer(cont)
	return vzAny, nil
}

func (sm *ServiceManager) DeleteAny(vrf string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ctx-%s/any", tenant, vrf)
	return sm.DeleteByDn(dn, models.VzanyClassName)
}

func (sm *ServiceManager) UpdateAny(vrf string, tenant string, description string, vzAnyattr models.AnyAttributes) (*models.Any, error) {
	rn := fmt.Sprintf("any")
	parentDn := fmt.Sprintf("uni/tn-%s/ctx-%s", tenant, vrf)
	vzAny := models.NewAny(rn, parentDn, description, vzAnyattr)

	vzAny.Status = "modified"
	err := sm.Save(vzAny)
	return vzAny, err

}

func (sm *ServiceManager) ListAny(vrf string, tenant string) ([]*models.Any, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctx-%s/vzAny.json", baseurlStr, tenant, vrf)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.AnyListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvzRsAnyToConsFromAny(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsanyToCons-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "vzRsAnyToCons", dn))

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

func (sm *ServiceManager) DeleteRelationvzRsAnyToConsFromAny(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsanyToCons-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "vzRsAnyToCons")
}

func (sm *ServiceManager) ReadRelationvzRsAnyToConsFromAny(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsAnyToCons")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsAnyToCons")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationvzRsAnyToConsIfFromAny(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsanyToConsIf-%s", parentDn, tnVzCPIfName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "vzRsAnyToConsIf", dn))

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

func (sm *ServiceManager) DeleteRelationvzRsAnyToConsIfFromAny(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsanyToConsIf-%s", parentDn, tnVzCPIfName)
	return sm.DeleteByDn(dn, "vzRsAnyToConsIf")
}

func (sm *ServiceManager) ReadRelationvzRsAnyToConsIfFromAny(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsAnyToConsIf")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "tDn")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationvzRsAnyToProvFromAny(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsanyToProv-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "vzRsAnyToProv", dn))

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

func (sm *ServiceManager) DeleteRelationvzRsAnyToProvFromAny(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsanyToProv-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "vzRsAnyToProv")
}

func (sm *ServiceManager) ReadRelationvzRsAnyToProvFromAny(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsAnyToProv")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsAnyToProv")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
