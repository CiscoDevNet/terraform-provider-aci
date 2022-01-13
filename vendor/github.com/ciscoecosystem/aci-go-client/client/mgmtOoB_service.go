package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateOutOfBandManagementEPg(name string, management_profile string, tenant string, description string, mgmtOoBattr models.OutOfBandManagementEPgAttributes) (*models.OutOfBandManagementEPg, error) {
	rn := fmt.Sprintf("oob-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/mgmtp-%s", tenant, management_profile)
	mgmtOoB := models.NewOutOfBandManagementEPg(rn, parentDn, description, mgmtOoBattr)
	err := sm.Save(mgmtOoB)
	return mgmtOoB, err
}

func (sm *ServiceManager) ReadOutOfBandManagementEPg(name string, management_profile string, tenant string) (*models.OutOfBandManagementEPg, error) {
	dn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/oob-%s", tenant, management_profile, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	mgmtOoB := models.OutOfBandManagementEPgFromContainer(cont)
	return mgmtOoB, nil
}

func (sm *ServiceManager) DeleteOutOfBandManagementEPg(name string, management_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/oob-%s", tenant, management_profile, name)
	return sm.DeleteByDn(dn, models.MgmtoobClassName)
}

func (sm *ServiceManager) UpdateOutOfBandManagementEPg(name string, management_profile string, tenant string, description string, mgmtOoBattr models.OutOfBandManagementEPgAttributes) (*models.OutOfBandManagementEPg, error) {
	rn := fmt.Sprintf("oob-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/mgmtp-%s", tenant, management_profile)
	mgmtOoB := models.NewOutOfBandManagementEPg(rn, parentDn, description, mgmtOoBattr)

	mgmtOoB.Status = "modified"
	err := sm.Save(mgmtOoB)
	return mgmtOoB, err

}

func (sm *ServiceManager) ListOutOfBandManagementEPg(management_profile string, tenant string) ([]*models.OutOfBandManagementEPg, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/mgmtp-%s/mgmtOoB.json", baseurlStr, tenant, management_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.OutOfBandManagementEPgListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationmgmtRsOoBProvFromOutOfBandManagementEPg(parentDn, tnVzOOBBrCPName string) error {
	dn := fmt.Sprintf("%s/rsooBProv-%s", parentDn, tnVzOOBBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzOOBBrCPName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "mgmtRsOoBProv", dn, tnVzOOBBrCPName))

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

func (sm *ServiceManager) DeleteRelationmgmtRsOoBProvFromOutOfBandManagementEPg(parentDn, tnVzOOBBrCPName string) error {
	dn := fmt.Sprintf("%s/rsooBProv-%s", parentDn, tnVzOOBBrCPName)
	return sm.DeleteByDn(dn, "mgmtRsOoBProv")
}

func (sm *ServiceManager) ReadRelationmgmtRsOoBProvFromOutOfBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "mgmtRsOoBProv")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "mgmtRsOoBProv")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationmgmtRsOoBStNodeFromOutOfBandManagementEPg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsooBStNode-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "mgmtRsOoBStNode", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationmgmtRsOoBStNodeFromOutOfBandManagementEPg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsooBStNode-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "mgmtRsOoBStNode")
}

func (sm *ServiceManager) ReadRelationmgmtRsOoBStNodeFromOutOfBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "mgmtRsOoBStNode")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "mgmtRsOoBStNode")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationmgmtRsOoBCtxFromOutOfBandManagementEPg(parentDn, tnFvCtxName string) error {
	dn := fmt.Sprintf("%s/rsooBCtx", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnFvCtxName": "%s"}
		}
	}`, "mgmtRsOoBCtx", dn, tnFvCtxName))

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

func (sm *ServiceManager) ReadRelationmgmtRsOoBCtxFromOutOfBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "mgmtRsOoBCtx")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "mgmtRsOoBCtx")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
