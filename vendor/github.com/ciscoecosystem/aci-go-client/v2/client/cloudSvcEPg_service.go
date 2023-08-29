package client

import (
	"encoding/json"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateCloudServiceEPg(name string, cloud_application_container string, tenant string, description string, cloudSvcEPgAttr models.CloudServiceEPgAttributes) (*models.CloudServiceEPg, error) {

	rn := fmt.Sprintf(models.RnCloudSvcEPg, name)

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPg, tenant, cloud_application_container)
	cloudSvcEPg := models.NewCloudServiceEPg(rn, parentDn, description, cloudSvcEPgAttr)

	err := sm.Save(cloudSvcEPg)
	return cloudSvcEPg, err
}

func (sm *ServiceManager) ReadCloudServiceEPg(name string, cloud_application_container string, tenant string) (*models.CloudServiceEPg, error) {

	rn := fmt.Sprintf(models.RnCloudSvcEPg, name)

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPg, tenant, cloud_application_container)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudSvcEPg := models.CloudServiceEPgFromContainer(cont)
	return cloudSvcEPg, nil
}

func (sm *ServiceManager) DeleteCloudServiceEPg(name string, cloud_application_container string, tenant string) error {

	rn := fmt.Sprintf(models.RnCloudSvcEPg, name)

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPg, tenant, cloud_application_container)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	return sm.DeleteByDn(dn, models.CloudSvcEPgClassName)
}

func (sm *ServiceManager) UpdateCloudServiceEPg(name string, cloud_application_container string, tenant string, description string, cloudSvcEPgAttr models.CloudServiceEPgAttributes) (*models.CloudServiceEPg, error) {

	rn := fmt.Sprintf(models.RnCloudSvcEPg, name)

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPg, tenant, cloud_application_container)
	cloudSvcEPg := models.NewCloudServiceEPg(rn, parentDn, description, cloudSvcEPgAttr)

	cloudSvcEPg.Status = "modified"
	err := sm.Save(cloudSvcEPg)
	return cloudSvcEPg, err
}

func (sm *ServiceManager) ListCloudServiceEPg(cloud_application_container string, tenant string) ([]*models.CloudServiceEPg, error) {

	parentDn := fmt.Sprintf(models.ParentDnCloudSvcEPg, tenant, cloud_application_container)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.CloudSvcEPgClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CloudServiceEPgListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationcloudRsCloudEPgCtxFromCloudServiceEpg(parentDn, annotation, tnFvCtxName string) error {
	dn := fmt.Sprintf("%s/rsCloudEPgCtx", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnFvCtxName": "%s"	
			}
		}
	}`, "cloudRsCloudEPgCtx", dn, annotation, tnFvCtxName))

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

func (sm *ServiceManager) DeleteRelationcloudRsCloudEPgCtxFromCloudServiceEpg(parentDn string) error {
	dn := fmt.Sprintf("%s/rsCloudEPgCtx", parentDn)
	return sm.DeleteByDn(dn, "cloudRsCloudEPgCtx")
}

func (sm *ServiceManager) ReadRelationcloudRsCloudEPgCtxFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "cloudRsCloudEPgCtx")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "cloudRsCloudEPgCtx")

	if len(contList) > 0 {
		paramMap := make(map[string]string)
		paramMap["tnFvCtxName"] = models.G(contList[0], "tnFvCtxName")
		paramMap["tDn"] = models.G(contList[0], "tDn")
		return paramMap, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationfvRsConsFromCloudServiceEpg(parentDn, annotation, prio string, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnVzBrCPName": "%s"	
			}
		}
	}`, "fvRsCons", dn, annotation, tnVzBrCPName))

	attributes := map[string]interface{}{
		"prio": prio,
	}
	var output map[string]interface{}
	err_output := json.Unmarshal([]byte(containerJSON), &output)
	if err_output != nil {
		return err_output
	}
	for _, mo := range output {
		if mo_map, ok := mo.(map[string]interface{}); ok {
			for _, mo_attributes := range mo_map {
				if mo_attributes_map, ok := mo_attributes.(map[string]interface{}); ok {
					for key, value := range attributes {
						if value != "" {
							mo_attributes_map[key] = value
						}
					}
				}
			}
		}

	}
	input, out_err := json.Marshal(output)
	if out_err != nil {
		return out_err
	}
	jsonPayload, err := container.ParseJSON(input)
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

func (sm *ServiceManager) DeleteRelationfvRsConsFromCloudServiceEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsCons")
}

func (sm *ServiceManager) ReadRelationfvRsConsFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsCons")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsCons")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnVzBrCPName"] = models.G(contItem, "tnVzBrCPName")
		paramMap["tDn"] = models.G(contList[0], "tDn")
		paramMap["prio"] = models.G(contList[0], "prio")
		st.Add(paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsConsIfFromCloudServiceEpg(parentDn, annotation, prio string, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnVzCPIfName": "%s"	
			}
		}
	}`, "fvRsConsIf", dn, annotation, tnVzCPIfName))

	attributes := map[string]interface{}{
		"prio": prio,
	}
	var output map[string]interface{}
	err_output := json.Unmarshal([]byte(containerJSON), &output)
	if err_output != nil {
		return err_output
	}
	for _, mo := range output {
		if mo_map, ok := mo.(map[string]interface{}); ok {
			for _, mo_attributes := range mo_map {
				if mo_attributes_map, ok := mo_attributes.(map[string]interface{}); ok {
					for key, value := range attributes {
						if value != "" {
							mo_attributes_map[key] = value
						}
					}
				}
			}
		}

	}
	input, out_err := json.Marshal(output)
	if out_err != nil {
		return out_err
	}
	jsonPayload, err := container.ParseJSON(input)
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

func (sm *ServiceManager) DeleteRelationfvRsConsIfFromCloudServiceEpg(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	return sm.DeleteByDn(dn, "fvRsConsIf")
}

func (sm *ServiceManager) ReadRelationfvRsConsIfFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsConsIf")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsConsIf")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnVzCPIfName"] = models.G(contItem, "tnVzCPIfName")
		paramMap["tDn"] = models.G(contList[0], "tDn")
		paramMap["prio"] = models.G(contList[0], "prio")
		st.Add(paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsCustQosPolFromCloudServiceEpg(parentDn, annotation, tnQosCustomPolName string) error {
	dn := fmt.Sprintf("%s/rscustQosPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnQosCustomPolName": "%s"	
			}
		}
	}`, "fvRsCustQosPol", dn, annotation, tnQosCustomPolName))

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

func (sm *ServiceManager) DeleteRelationfvRsCustQosPolFromCloudServiceEpg(parentDn string) error {
	dn := fmt.Sprintf("%s/rscustQosPol", parentDn)
	return sm.DeleteByDn(dn, "fvRsCustQosPol")
}

func (sm *ServiceManager) ReadRelationfvRsCustQosPolFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsCustQosPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsCustQosPol")

	if len(contList) > 0 {
		paramMap := make(map[string]string)
		paramMap["tnQosCustomPolName"] = models.G(contList[0], "tnQosCustomPolName")
		paramMap["tDn"] = models.G(contList[0], "tDn")
		return paramMap, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationfvRsGraphDefFromCloudServiceEpg(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsgraphDef-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"	
			}
		}
	}`, "fvRsGraphDef", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationfvRsGraphDefFromCloudServiceEpg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsgraphDef-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsGraphDef")
}

func (sm *ServiceManager) ReadRelationfvRsGraphDefFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsGraphDef")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsGraphDef")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tDn"] = models.G(contItem, "tDn")
		st.Add(paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsIntraEpgFromCloudServiceEpg(parentDn, annotation, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnVzBrCPName": "%s"	
			}
		}
	}`, "fvRsIntraEpg", dn, annotation, tnVzBrCPName))

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

func (sm *ServiceManager) DeleteRelationfvRsIntraEpgFromCloudServiceEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsIntraEpg")
}

func (sm *ServiceManager) ReadRelationfvRsIntraEpgFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsIntraEpg")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsIntraEpg")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnVzBrCPName"] = models.G(contItem, "tnVzBrCPName")
		paramMap["tDn"] = models.G(contList[0], "tDn")
		st.Add(paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsProtByFromCloudServiceEpg(parentDn, annotation, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnVzTabooName": "%s"	
			}
		}
	}`, "fvRsProtBy", dn, annotation, tnVzTabooName))

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

func (sm *ServiceManager) DeleteRelationfvRsProtByFromCloudServiceEpg(parentDn, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	return sm.DeleteByDn(dn, "fvRsProtBy")
}

func (sm *ServiceManager) ReadRelationfvRsProtByFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsProtBy")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsProtBy")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnVzTabooName"] = models.G(contItem, "tnVzTabooName")
		paramMap["tDn"] = models.G(contList[0], "tDn")
		st.Add(paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsProvFromCloudServiceEpg(parentDn, annotation, matchT string, prio string, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnVzBrCPName": "%s"	
			}
		}
	}`, "fvRsProv", dn, annotation, tnVzBrCPName))

	attributes := map[string]interface{}{
		"matchT": matchT,
		"prio":   prio,
	}
	var output map[string]interface{}
	err_output := json.Unmarshal([]byte(containerJSON), &output)
	if err_output != nil {
		return err_output
	}
	for _, mo := range output {
		if mo_map, ok := mo.(map[string]interface{}); ok {
			for _, mo_attributes := range mo_map {
				if mo_attributes_map, ok := mo_attributes.(map[string]interface{}); ok {
					for key, value := range attributes {
						if value != "" {
							mo_attributes_map[key] = value
						}
					}
				}
			}
		}

	}
	input, out_err := json.Marshal(output)
	if out_err != nil {
		return out_err
	}
	jsonPayload, err := container.ParseJSON(input)
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

func (sm *ServiceManager) DeleteRelationfvRsProvFromCloudServiceEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsProv")
}

func (sm *ServiceManager) ReadRelationfvRsProvFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsProv")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsProv")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnVzBrCPName"] = models.G(contItem, "tnVzBrCPName")
		paramMap["tDn"] = models.G(contList[0], "tDn")
		paramMap["matchT"] = models.G(contList[0], "matchT")
		paramMap["prio"] = models.G(contList[0], "prio")
		st.Add(paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsProvDefFromCloudServiceEpg(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsprovDef-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"	
			}
		}
	}`, "fvRsProvDef", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationfvRsProvDefFromCloudServiceEpg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsprovDef-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsProvDef")
}

func (sm *ServiceManager) ReadRelationfvRsProvDefFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsProvDef")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsProvDef")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tDn"] = models.G(contItem, "tDn")
		st.Add(paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsSecInheritedFromCloudServiceEpg(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"	
			}
		}
	}`, "fvRsSecInherited", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationfvRsSecInheritedFromCloudServiceEpg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsSecInherited")
}

func (sm *ServiceManager) ReadRelationfvRsSecInheritedFromCloudServiceEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsSecInherited")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsSecInherited")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tDn"] = models.G(contItem, "tDn")
		st.Add(paramMap)
	}
	return st, err
}
