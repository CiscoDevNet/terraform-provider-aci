package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateExternalNetworkInstanceProfile(name string, l3_outside string, tenant string, description string, l3extInstPattr models.ExternalNetworkInstanceProfileAttributes) (*models.ExternalNetworkInstanceProfile, error) {
	rn := fmt.Sprintf("instP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s", tenant, l3_outside)
	l3extInstP := models.NewExternalNetworkInstanceProfile(rn, parentDn, description, l3extInstPattr)
	err := sm.Save(l3extInstP)
	return l3extInstP, err
}

func (sm *ServiceManager) ReadExternalNetworkInstanceProfile(name string, l3_outside string, tenant string) (*models.ExternalNetworkInstanceProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/instP-%s", tenant, l3_outside, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extInstP := models.ExternalNetworkInstanceProfileFromContainer(cont)
	return l3extInstP, nil
}

func (sm *ServiceManager) DeleteExternalNetworkInstanceProfile(name string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/instP-%s", tenant, l3_outside, name)
	return sm.DeleteByDn(dn, models.L3extinstpClassName)
}

func (sm *ServiceManager) UpdateExternalNetworkInstanceProfile(name string, l3_outside string, tenant string, description string, l3extInstPattr models.ExternalNetworkInstanceProfileAttributes) (*models.ExternalNetworkInstanceProfile, error) {
	rn := fmt.Sprintf("instP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s", tenant, l3_outside)
	l3extInstP := models.NewExternalNetworkInstanceProfile(rn, parentDn, description, l3extInstPattr)

	l3extInstP.Status = "modified"
	err := sm.Save(l3extInstP)
	return l3extInstP, err

}

func (sm *ServiceManager) ListExternalNetworkInstanceProfile(l3_outside string, tenant string) ([]*models.ExternalNetworkInstanceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/l3extInstP.json", baseurlStr, tenant, l3_outside)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ExternalNetworkInstanceProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsSecInherited", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsSecInherited")
}

func (sm *ServiceManager) ReadRelationfvRsSecInheritedFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationfvRsProvFromExternalNetworkInstanceProfile(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsProv", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsProvFromExternalNetworkInstanceProfile(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsProv")
}

func (sm *ServiceManager) ReadRelationfvRsProvFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(parentDn, tnExtnwDomPName string) error {
	dn := fmt.Sprintf("%s/rsl3InstPToDomP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s"
								
			}
		}
	}`, "l3extRsL3InstPToDomP", dn, tnExtnwDomPName))

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

func (sm *ServiceManager) ReadRelationl3extRsL3InstPToDomPFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsL3InstPToDomP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsL3InstPToDomP")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(parentDn, tnFvAEPgName string) error {
	dn := fmt.Sprintf("%s/rsInstPToNatMappingEPg", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsInstPToNatMappingEPg", dn, tnFvAEPgName))

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

func (sm *ServiceManager) DeleteRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(parentDn string) error {
	dn := fmt.Sprintf("%s/rsInstPToNatMappingEPg", parentDn)
	return sm.DeleteByDn(dn, "l3extRsInstPToNatMappingEPg")
}

func (sm *ServiceManager) ReadRelationl3extRsInstPToNatMappingEPgFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsInstPToNatMappingEPg")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsInstPToNatMappingEPg")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsConsIfFromExternalNetworkInstanceProfile(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsConsIf", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsConsIfFromExternalNetworkInstanceProfile(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	return sm.DeleteByDn(dn, "fvRsConsIf")
}

func (sm *ServiceManager) ReadRelationfvRsConsIfFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(parentDn, tnQosCustomPolName string) error {
	dn := fmt.Sprintf("%s/rscustQosPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosCustomPolName": "%s","annotation":"orchestrator:terraform"
								
			}
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) ReadRelationfvRsCustQosPolFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(parentDn, tnRtctrlProfileName, direction string) error {
	dn := fmt.Sprintf("%s/rsinstPToProfile-[%s]-%s", parentDn, tnRtctrlProfileName, direction)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "l3extRsInstPToProfile", dn))

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

func (sm *ServiceManager) DeleteRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(parentDn, tnRtctrlProfileName, direction string) error {
	dn := fmt.Sprintf("%s/rsinstPToProfile-[%s]-%s", parentDn, tnRtctrlProfileName, direction)
	return sm.DeleteByDn(dn, "l3extRsInstPToProfile")
}

func (sm *ServiceManager) ReadRelationl3extRsInstPToProfileFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsInstPToProfile")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsInstPToProfile")

	st := make([]map[string]string, 0)

	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnRtctrlProfileName"] = models.G(contItem, "tDn")
		paramMap["direction"] = models.G(contItem, "direction")

		st = append(st, paramMap)

	}

	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsConsFromExternalNetworkInstanceProfile(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsCons", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsConsFromExternalNetworkInstanceProfile(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsCons")
}

func (sm *ServiceManager) ReadRelationfvRsConsFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationfvRsProtByFromExternalNetworkInstanceProfile(parentDn, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsProtBy", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsProtByFromExternalNetworkInstanceProfile(parentDn, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	return sm.DeleteByDn(dn, "fvRsProtBy")
}

func (sm *ServiceManager) ReadRelationfvRsProtByFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsIntraEpg", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsIntraEpg")
}

func (sm *ServiceManager) ReadRelationfvRsIntraEpgFromExternalNetworkInstanceProfile(parentDn string) (interface{}, error) {
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
