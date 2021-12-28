package client

import (
	"encoding/json"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateVSwitchPolicyGroup(vmm_domain string, provider_profile_vendor string, description string, nameAlias string, vmmVSwitchPolicyContAttr models.VSwitchPolicyGroupAttributes) (*models.VSwitchPolicyGroup, error) {
	rn := fmt.Sprintf(models.RnvmmVSwitchPolicyCont)
	parentDn := fmt.Sprintf(models.ParentDnvmmVSwitchPolicyCont, provider_profile_vendor, vmm_domain)
	vmmVSwitchPolicyCont := models.NewVSwitchPolicyGroup(rn, parentDn, description, nameAlias, vmmVSwitchPolicyContAttr)
	err := sm.Save(vmmVSwitchPolicyCont)
	return vmmVSwitchPolicyCont, err
}

func (sm *ServiceManager) ReadVSwitchPolicyGroup(vmm_domain string, provider_profile_vendor string) (*models.VSwitchPolicyGroup, error) {
	dn := fmt.Sprintf(models.DnvmmVSwitchPolicyCont, provider_profile_vendor, vmm_domain)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	vmmVSwitchPolicyCont := models.VSwitchPolicyGroupFromContainer(cont)
	return vmmVSwitchPolicyCont, nil
}

func (sm *ServiceManager) DeleteVSwitchPolicyGroup(vmm_domain string, provider_profile_vendor string) error {
	dn := fmt.Sprintf(models.DnvmmVSwitchPolicyCont, provider_profile_vendor, vmm_domain)
	return sm.DeleteByDn(dn, models.VmmvswitchpolicycontClassName)
}

func (sm *ServiceManager) UpdateVSwitchPolicyGroup(vmm_domain string, provider_profile_vendor string, description string, nameAlias string, vmmVSwitchPolicyContAttr models.VSwitchPolicyGroupAttributes) (*models.VSwitchPolicyGroup, error) {
	rn := fmt.Sprintf(models.RnvmmVSwitchPolicyCont)
	parentDn := fmt.Sprintf(models.ParentDnvmmVSwitchPolicyCont, provider_profile_vendor, vmm_domain)
	vmmVSwitchPolicyCont := models.NewVSwitchPolicyGroup(rn, parentDn, description, nameAlias, vmmVSwitchPolicyContAttr)
	vmmVSwitchPolicyCont.Status = "modified"
	err := sm.Save(vmmVSwitchPolicyCont)
	return vmmVSwitchPolicyCont, err
}

func (sm *ServiceManager) ListVSwitchPolicyGroup(vmm_domain string, provider_profile_vendor string) ([]*models.VSwitchPolicyGroup, error) {
	dnUrl := fmt.Sprintf("%s/uni/vmmp-%s/dom-%s/vmmVSwitchPolicyCont.json", models.BaseurlStr, provider_profile_vendor, vmm_domain)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.VSwitchPolicyGroupListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvmmRsVswitchExporterPol(parentDn, annotation, activeFlowTimeOut string, idleFlowTimeOut string, samplingRate string, tDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchExporterPol-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVswitchExporterPol", dn, annotation, tDn))

	attributes := map[string]interface{}{
		"activeFlowTimeOut": activeFlowTimeOut,
		"idleFlowTimeOut":   idleFlowTimeOut,
		"samplingRate":      samplingRate,
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

func (sm *ServiceManager) DeleteRelationvmmRsVswitchExporterPol(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchExporterPol-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "vmmRsVswitchExporterPol")
}

func (sm *ServiceManager) ReadRelationvmmRsVswitchExporterPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVswitchExporterPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVswitchExporterPol")

	st := make([]map[string]string, 0, 1)
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["activeFlowTimeOut"] = models.G(contItem, "activeFlowTimeOut")
		paramMap["idleFlowTimeOut"] = models.G(contItem, "idleFlowTimeOut")
		paramMap["samplingRate"] = models.G(contItem, "samplingRate")
		paramMap["tDn"] = models.G(contItem, "tDn")
		st = append(st, paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationvmmRsVswitchOverrideCdpIfPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideCdpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVswitchOverrideCdpIfPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsVswitchOverrideCdpIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideCdpIfPol", parentDn)
	return sm.DeleteByDn(dn, "vmmRsVswitchOverrideCdpIfPol")
}

func (sm *ServiceManager) ReadRelationvmmRsVswitchOverrideCdpIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVswitchOverrideCdpIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVswitchOverrideCdpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsVswitchOverrideFwPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideFwPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVswitchOverrideFwPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsVswitchOverrideFwPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideFwPol", parentDn)
	return sm.DeleteByDn(dn, "vmmRsVswitchOverrideFwPol")
}

func (sm *ServiceManager) ReadRelationvmmRsVswitchOverrideFwPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVswitchOverrideFwPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVswitchOverrideFwPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsVswitchOverrideLacpPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideLacpPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVswitchOverrideLacpPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsVswitchOverrideLacpPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideLacpPol", parentDn)
	return sm.DeleteByDn(dn, "vmmRsVswitchOverrideLacpPol")
}

func (sm *ServiceManager) ReadRelationvmmRsVswitchOverrideLacpPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVswitchOverrideLacpPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVswitchOverrideLacpPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsVswitchOverrideLldpIfPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideLldpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVswitchOverrideLldpIfPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsVswitchOverrideLldpIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideLldpIfPol", parentDn)
	return sm.DeleteByDn(dn, "vmmRsVswitchOverrideLldpIfPol")
}

func (sm *ServiceManager) ReadRelationvmmRsVswitchOverrideLldpIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVswitchOverrideLldpIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVswitchOverrideLldpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsVswitchOverrideMcpIfPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideMcpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVswitchOverrideMcpIfPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsVswitchOverrideMcpIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideMcpIfPol", parentDn)
	return sm.DeleteByDn(dn, "vmmRsVswitchOverrideMcpIfPol")
}

func (sm *ServiceManager) ReadRelationvmmRsVswitchOverrideMcpIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVswitchOverrideMcpIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVswitchOverrideMcpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsVswitchOverrideMtuPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideMtuPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVswitchOverrideMtuPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsVswitchOverrideMtuPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideMtuPol", parentDn)
	return sm.DeleteByDn(dn, "vmmRsVswitchOverrideMtuPol")
}

func (sm *ServiceManager) ReadRelationvmmRsVswitchOverrideMtuPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVswitchOverrideMtuPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVswitchOverrideMtuPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvmmRsVswitchOverrideStpPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideStpPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vmmRsVswitchOverrideStpPol", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvmmRsVswitchOverrideStpPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvswitchOverrideStpPol", parentDn)
	return sm.DeleteByDn(dn, "vmmRsVswitchOverrideStpPol")
}

func (sm *ServiceManager) ReadRelationvmmRsVswitchOverrideStpPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vmmRsVswitchOverrideStpPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vmmRsVswitchOverrideStpPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
