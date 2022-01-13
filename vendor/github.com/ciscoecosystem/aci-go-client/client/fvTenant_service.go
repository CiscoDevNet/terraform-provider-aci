package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateTenant(name string, description string, fvTenantattr models.TenantAttributes) (*models.Tenant, error) {
	rn := fmt.Sprintf("tn-%s", name)
	parentDn := fmt.Sprintf("uni")
	fvTenant := models.NewTenant(rn, parentDn, description, fvTenantattr)
	err := sm.Save(fvTenant)
	return fvTenant, err
}

func (sm *ServiceManager) ReadTenant(name string) (*models.Tenant, error) {
	dn := fmt.Sprintf("uni/tn-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvTenant := models.TenantFromContainer(cont)
	return fvTenant, nil
}

func (sm *ServiceManager) DeleteTenant(name string) error {
	dn := fmt.Sprintf("uni/tn-%s", name)
	return sm.DeleteByDn(dn, models.FvtenantClassName)
}

func (sm *ServiceManager) UpdateTenant(name string, description string, fvTenantattr models.TenantAttributes) (*models.Tenant, error) {
	rn := fmt.Sprintf("tn-%s", name)
	parentDn := fmt.Sprintf("uni")
	fvTenant := models.NewTenant(rn, parentDn, description, fvTenantattr)

	fvTenant.Status = "modified"
	err := sm.Save(fvTenant)
	return fvTenant, err

}

func (sm *ServiceManager) ListTenant() ([]*models.Tenant, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fvTenant.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.TenantListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationfvRsTnDenyRuleFromTenant(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rstnDenyRule-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsTnDenyRule", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsTnDenyRuleFromTenant(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rstnDenyRule-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsTnDenyRule")
}

func (sm *ServiceManager) ReadRelationfvRsTnDenyRuleFromTenant(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsTnDenyRule")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsTnDenyRule")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsTenantMonPolFromTenant(parentDn, tnMonEPGPolName string) error {
	dn := fmt.Sprintf("%s/rsTenantMonPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMonEPGPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsTenantMonPol", dn, tnMonEPGPolName))

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

func (sm *ServiceManager) ReadRelationfvRsTenantMonPolFromTenant(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsTenantMonPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsTenantMonPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
