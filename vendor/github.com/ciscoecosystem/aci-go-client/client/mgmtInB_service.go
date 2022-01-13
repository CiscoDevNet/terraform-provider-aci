package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateInBandManagementEPg(name string, management_profile string, tenant string, description string, mgmtInBattr models.InBandManagementEPgAttributes) (*models.InBandManagementEPg, error) {
	rn := fmt.Sprintf("inb-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/mgmtp-%s", tenant, management_profile)
	mgmtInB := models.NewInBandManagementEPg(rn, parentDn, description, mgmtInBattr)
	err := sm.Save(mgmtInB)
	return mgmtInB, err
}

func (sm *ServiceManager) ReadInBandManagementEPg(name string, management_profile string, tenant string) (*models.InBandManagementEPg, error) {
	dn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/inb-%s", tenant, management_profile, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	mgmtInB := models.InBandManagementEPgFromContainer(cont)
	return mgmtInB, nil
}

func (sm *ServiceManager) DeleteInBandManagementEPg(name string, management_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/mgmtp-%s/inb-%s", tenant, management_profile, name)
	return sm.DeleteByDn(dn, models.MgmtinbClassName)
}

func (sm *ServiceManager) UpdateInBandManagementEPg(name string, management_profile string, tenant string, description string, mgmtInBattr models.InBandManagementEPgAttributes) (*models.InBandManagementEPg, error) {
	rn := fmt.Sprintf("inb-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/mgmtp-%s", tenant, management_profile)
	mgmtInB := models.NewInBandManagementEPg(rn, parentDn, description, mgmtInBattr)

	mgmtInB.Status = "modified"
	err := sm.Save(mgmtInB)
	return mgmtInB, err

}

func (sm *ServiceManager) ListInBandManagementEPg(management_profile string, tenant string) ([]*models.InBandManagementEPg, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/mgmtp-%s/mgmtInB.json", baseurlStr, tenant, management_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.InBandManagementEPgListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationfvRsSecInheritedFromInBandManagementEPg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "fvRsSecInherited", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationfvRsSecInheritedFromInBandManagementEPg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsSecInherited")
}

func (sm *ServiceManager) ReadRelationfvRsSecInheritedFromInBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsSecInherited")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsSecInherited")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsProvFromInBandManagementEPg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzBrCPName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "fvRsProv", dn, tnVzBrCPName))

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

func (sm *ServiceManager) DeleteRelationfvRsProvFromInBandManagementEPg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsProv")
}

func (sm *ServiceManager) ReadRelationfvRsProvFromInBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsProv")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsProv")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsConsIfFromInBandManagementEPg(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzCPIfName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "fvRsConsIf", dn, tnVzCPIfName))

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

func (sm *ServiceManager) DeleteRelationfvRsConsIfFromInBandManagementEPg(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	return sm.DeleteByDn(dn, "fvRsConsIf")
}

func (sm *ServiceManager) ReadRelationfvRsConsIfFromInBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsConsIf")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsConsIf")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsCustQosPolFromInBandManagementEPg(parentDn, tnQosCustomPolName string) error {
	dn := fmt.Sprintf("%s/rscustQosPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnQosCustomPolName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "fvRsCustQosPol", dn, tnQosCustomPolName))

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

func (sm *ServiceManager) ReadRelationfvRsCustQosPolFromInBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsCustQosPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsCustQosPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationmgmtRsMgmtBDFromInBandManagementEPg(parentDn, tnFvBDName string) error {
	dn := fmt.Sprintf("%s/rsmgmtBD", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnFvBDName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "mgmtRsMgmtBD", dn, tnFvBDName))

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

func (sm *ServiceManager) ReadRelationmgmtRsMgmtBDFromInBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "mgmtRsMgmtBD")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "mgmtRsMgmtBD")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsConsFromInBandManagementEPg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzBrCPName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "fvRsCons", dn, tnVzBrCPName))

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

func (sm *ServiceManager) DeleteRelationfvRsConsFromInBandManagementEPg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsCons")
}

func (sm *ServiceManager) ReadRelationfvRsConsFromInBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsCons")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsCons")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsProtByFromInBandManagementEPg(parentDn, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzTabooName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "fvRsProtBy", dn, tnVzTabooName))

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

func (sm *ServiceManager) DeleteRelationfvRsProtByFromInBandManagementEPg(parentDn, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	return sm.DeleteByDn(dn, "fvRsProtBy")
}

func (sm *ServiceManager) ReadRelationfvRsProtByFromInBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsProtBy")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsProtBy")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationmgmtRsInBStNodeFromInBandManagementEPg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsinBStNode-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "mgmtRsInBStNode", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationmgmtRsInBStNodeFromInBandManagementEPg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsinBStNode-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "mgmtRsInBStNode")
}

func (sm *ServiceManager) ReadRelationmgmtRsInBStNodeFromInBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "mgmtRsInBStNode")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "mgmtRsInBStNode")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsIntraEpgFromInBandManagementEPg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzBrCPName": "%s", 
				"annotation":"orchestrator:terraform"}
		}
	}`, "fvRsIntraEpg", dn, tnVzBrCPName))

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

func (sm *ServiceManager) DeleteRelationfvRsIntraEpgFromInBandManagementEPg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsIntraEpg")
}

func (sm *ServiceManager) ReadRelationfvRsIntraEpgFromInBandManagementEPg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsIntraEpg")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsIntraEpg")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
