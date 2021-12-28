package client

import (
	"encoding/json"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateVMMController(name string, vmm_domain string, provider_profile_vendor string, description string, nameAlias string, vmmCtrlrPAttr models.VMMControllerAttributes) (*models.VMMController, error) {
	rn := fmt.Sprintf(models.RnvmmCtrlrP, name)
	parentDn := fmt.Sprintf(models.ParentDnvmmCtrlrP, provider_profile_vendor, vmm_domain)
	vmmCtrlrP := models.NewVMMController(rn, parentDn, nameAlias, vmmCtrlrPAttr)
	err := sm.Save(vmmCtrlrP)
	return vmmCtrlrP, err
}

func (sm *ServiceManager) ReadVMMController(name string, vmm_domain string, provider_profile_vendor string) (*models.VMMController, error) {
	dn := fmt.Sprintf(models.DnvmmCtrlrP, provider_profile_vendor, vmm_domain, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	vmmCtrlrP := models.VMMControllerFromContainer(cont)
	return vmmCtrlrP, nil
}

func (sm *ServiceManager) DeleteVMMController(name string, vmm_domain string, provider_profile_vendor string) error {
	dn := fmt.Sprintf(models.DnvmmCtrlrP, provider_profile_vendor, vmm_domain, name)
	return sm.DeleteByDn(dn, models.VmmctrlrpClassName)
}

func (sm *ServiceManager) UpdateVMMController(name string, vmm_domain string, provider_profile_vendor string, description string, nameAlias string, vmmCtrlrPAttr models.VMMControllerAttributes) (*models.VMMController, error) {
	rn := fmt.Sprintf(models.RnvmmCtrlrP, name)
	parentDn := fmt.Sprintf(models.ParentDnvmmCtrlrP, provider_profile_vendor, vmm_domain)
	vmmCtrlrP := models.NewVMMController(rn, parentDn, nameAlias, vmmCtrlrPAttr)
	vmmCtrlrP.Status = "modified"
	err := sm.Save(vmmCtrlrP)
	return vmmCtrlrP, err
}

func (sm *ServiceManager) ListVMMController(vmm_domain string, provider_profile_vendor string) ([]*models.VMMController, error) {
	dnUrl := fmt.Sprintf("%s/uni/vmmp-%s/dom-%s/vmmCtrlrP.json", models.BaseurlStr, provider_profile_vendor, vmm_domain)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.VMMControllerListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvmmRsAcc(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsacc", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsAcc", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsAcc(parentDn string) error {
	dn := fmt.Sprintf("%s/rsacc", parentDn)
	return sm.DeleteByDn(dn, "vmmRsAcc")
}

func (sm *ServiceManager) ReadRelationvmmRsAcc(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsAcc")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsAcc")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsCtrlrPMonPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsctrlrPMonPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsCtrlrPMonPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsCtrlrPMonPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsctrlrPMonPol", parentDn)
	return sm.DeleteByDn(dn, "vmmRsCtrlrPMonPol")
}

func (sm *ServiceManager) ReadRelationvmmRsCtrlrPMonPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsCtrlrPMonPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsCtrlrPMonPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsMcastAddrNs(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsmcastAddrNs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsMcastAddrNs", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsMcastAddrNs(parentDn string) error {
	dn := fmt.Sprintf("%s/rsmcastAddrNs", parentDn)
	return sm.DeleteByDn(dn, "vmmRsMcastAddrNs")
}

func (sm *ServiceManager) ReadRelationvmmRsMcastAddrNs(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsMcastAddrNs")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsMcastAddrNs")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsMgmtEPg(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsmgmtEPg", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsMgmtEPg", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsMgmtEPg(parentDn string) error {
	dn := fmt.Sprintf("%s/rsmgmtEPg", parentDn)
	return sm.DeleteByDn(dn, "vmmRsMgmtEPg")
}

func (sm *ServiceManager) ReadRelationvmmRsMgmtEPg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsMgmtEPg")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsMgmtEPg")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsToExtDevMgr(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rstoExtDevMgr-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsToExtDevMgr", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsToExtDevMgr(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rstoExtDevMgr-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "vmmRsToExtDevMgr")
}

func (sm *ServiceManager) ReadRelationvmmRsToExtDevMgr(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsToExtDevMgr")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsToExtDevMgr")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationvmmRsVmmCtrlrP(parentDn, annotation, epgDeplPref string, tDn string) error {
	dn := fmt.Sprintf("%s/rsvmmCtrlrP-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVmmCtrlrP", dn, annotation, tDn))

	attributes := map[string]interface{}{
		"epgDeplPref": epgDeplPref,
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

func (sm *ServiceManager) DeleteRelationvmmRsVmmCtrlrP(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsvmmCtrlrP-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "vmmRsVmmCtrlrP")
}

func (sm *ServiceManager) ReadRelationvmmRsVmmCtrlrP(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVmmCtrlrP")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVmmCtrlrP")

	st := make([]map[string]string, 0, 1)
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["epgDeplPref"] = models.G(contItem, "epgDeplPref")
		paramMap["tDn"] = models.G(contItem, "tDn")
		st = append(st, paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationvmmRsVxlanNs(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvxlanNs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVxlanNs", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsVxlanNs(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvxlanNs", parentDn)
	return sm.DeleteByDn(dn, "vmmRsVxlanNs")
}

func (sm *ServiceManager) ReadRelationvmmRsVxlanNs(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVxlanNs")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVxlanNs")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsVxlanNsDef(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvxlanNsDef", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVxlanNsDef", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsVxlanNsDef(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvxlanNsDef", parentDn)
	return sm.DeleteByDn(dn, "vmmRsVxlanNsDef")
}

func (sm *ServiceManager) ReadRelationvmmRsVxlanNsDef(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVxlanNsDef")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVxlanNsDef")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
