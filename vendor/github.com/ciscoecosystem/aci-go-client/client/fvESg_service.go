package client

import (
	"encoding/json"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateEndpointSecurityGroup(name string, application_profile string, tenant string, description string, nameAlias string, fvESgAttr models.EndpointSecurityGroupAttributes) (*models.EndpointSecurityGroup, error) {
	rn := fmt.Sprintf(models.RnfvESg, name)
	parentDn := fmt.Sprintf(models.ParentDnfvESg, tenant, application_profile)
	fvESg := models.NewEndpointSecurityGroup(rn, parentDn, description, nameAlias, fvESgAttr)
	err := sm.Save(fvESg)
	return fvESg, err
}

func (sm *ServiceManager) ReadEndpointSecurityGroup(name string, application_profile string, tenant string) (*models.EndpointSecurityGroup, error) {
	dn := fmt.Sprintf(models.DnfvESg, tenant, application_profile, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	fvESg := models.EndpointSecurityGroupFromContainer(cont)
	return fvESg, nil
}

func (sm *ServiceManager) DeleteEndpointSecurityGroup(name string, application_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnfvESg, tenant, application_profile, name)
	return sm.DeleteByDn(dn, models.FvesgClassName)
}

func (sm *ServiceManager) UpdateEndpointSecurityGroup(name string, application_profile string, tenant string, description string, nameAlias string, fvESgAttr models.EndpointSecurityGroupAttributes) (*models.EndpointSecurityGroup, error) {
	rn := fmt.Sprintf(models.RnfvESg, name)
	parentDn := fmt.Sprintf(models.ParentDnfvESg, tenant, application_profile)
	fvESg := models.NewEndpointSecurityGroup(rn, parentDn, description, nameAlias, fvESgAttr)
	fvESg.Status = "modified"
	err := sm.Save(fvESg)
	return fvESg, err
}

func (sm *ServiceManager) ListEndpointSecurityGroup(application_profile string, tenant string) ([]*models.EndpointSecurityGroup, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/fvESg.json", models.BaseurlStr, tenant, application_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.EndpointSecurityGroupListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationfvRsCons(parentDn, annotation, prio string, tnVzBrCPName string) error {
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
	for _, value := range output {
		if rec, ok := value.(map[string]interface{}); ok {
			for _, val2 := range rec {
				if rec2, ok := val2.(map[string]interface{}); ok {
					for key, value := range attributes {
						if value != "" {
							rec2[key] = value
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

func (sm *ServiceManager) DeleteRelationfvRsCons(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsCons")
}

func (sm *ServiceManager) ReadRelationfvRsCons(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsCons")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsCons")

	st := make([]map[string]string, 0, 1)
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["prio"] = models.G(contItem, "prio")
		paramMap["tDn"] = models.G(contItem, "tDn")
		st = append(st, paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsConsIf(parentDn, annotation, prio string, tnVzCPIfName string) error {
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
	for _, value := range output {
		if rec, ok := value.(map[string]interface{}); ok {
			for _, val2 := range rec {
				if rec2, ok := val2.(map[string]interface{}); ok {
					for key, value := range attributes {
						if value != "" {
							rec2[key] = value
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

func (sm *ServiceManager) DeleteRelationfvRsConsIf(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	return sm.DeleteByDn(dn, "fvRsConsIf")
}

func (sm *ServiceManager) ReadRelationfvRsConsIf(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsConsIf")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsConsIf")

	st := make([]map[string]string, 0, 1)
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["prio"] = models.G(contItem, "prio")
		paramMap["tDn"] = models.G(contItem, "tDn")
		st = append(st, paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsCustQosPol(parentDn, annotation, tnQosCustomPolName string) error {
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

func (sm *ServiceManager) DeleteRelationfvRsCustQosPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rscustQosPol", parentDn)
	return sm.DeleteByDn(dn, "fvRsCustQosPol")
}

func (sm *ServiceManager) ReadRelationfvRsCustQosPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsCustQosPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsCustQosPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationfvRsIntraEpg(parentDn, annotation, tnVzBrCPName string) error {
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

func (sm *ServiceManager) DeleteRelationfvRsIntraEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsIntraEpg")
}

func (sm *ServiceManager) ReadRelationfvRsIntraEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsIntraEpg")
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

func (sm *ServiceManager) CreateRelationfvRsProtBy(parentDn, annotation, tnVzTabooName string) error {
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

func (sm *ServiceManager) DeleteRelationfvRsProtBy(parentDn, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	return sm.DeleteByDn(dn, "fvRsProtBy")
}

func (sm *ServiceManager) ReadRelationfvRsProtBy(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsProtBy")
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

func (sm *ServiceManager) CreateRelationfvRsProv(parentDn, annotation, matchT string, prio string, tnVzBrCPName string) error {
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
	for _, value := range output {
		if rec, ok := value.(map[string]interface{}); ok {
			for _, val2 := range rec {
				if rec2, ok := val2.(map[string]interface{}); ok {
					for key, value := range attributes {
						if value != "" {
							rec2[key] = value
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

func (sm *ServiceManager) DeleteRelationfvRsProv(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsProv")
}

func (sm *ServiceManager) ReadRelationfvRsProv(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsProv")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsProv")

	st := make([]map[string]string, 0, 1)
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["prio"] = models.G(contItem, "prio")
		paramMap["tDn"] = models.G(contItem, "tDn")
		paramMap["matchT"] = models.G(contItem, "matchT")
		st = append(st, paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationfvRsScope(parentDn, annotation, tnFvCtxName string) error {
	dn := fmt.Sprintf("%s/rsscope", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnFvCtxName": "%s"
			}
		}
	}`, "fvRsScope", dn, annotation, tnFvCtxName))

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

func (sm *ServiceManager) DeleteRelationfvRsScope(parentDn string) error {
	dn := fmt.Sprintf("%s/rsscope", parentDn)
	return sm.DeleteByDn(dn, "fvRsScope")
}

func (sm *ServiceManager) ReadRelationfvRsScope(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsScope")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsScope")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationfvRsSecInherited(parentDn, annotation, tDn string) error {
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

func (sm *ServiceManager) DeleteRelationfvRsSecInherited(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsSecInherited")
}

func (sm *ServiceManager) ReadRelationfvRsSecInherited(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsSecInherited")
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
