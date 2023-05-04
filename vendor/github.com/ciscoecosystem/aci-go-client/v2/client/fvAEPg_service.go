package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateApplicationEPG(name string, application_profile string, tenant string, description string, fvAEPgattr models.ApplicationEPGAttributes) (*models.ApplicationEPG, error) {
	rn := fmt.Sprintf("epg-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s", tenant, application_profile)
	fvAEPg := models.NewApplicationEPG(rn, parentDn, description, fvAEPgattr)
	err := sm.Save(fvAEPg)
	return fvAEPg, err
}

func (sm *ServiceManager) ReadApplicationEPG(name string, application_profile string, tenant string) (*models.ApplicationEPG, error) {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, application_profile, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvAEPg := models.ApplicationEPGFromContainer(cont)
	return fvAEPg, nil
}

func (sm *ServiceManager) DeleteApplicationEPG(name string, application_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, application_profile, name)
	return sm.DeleteByDn(dn, models.FvaepgClassName)
}

func (sm *ServiceManager) UpdateApplicationEPG(name string, application_profile string, tenant string, description string, fvAEPgattr models.ApplicationEPGAttributes) (*models.ApplicationEPG, error) {
	rn := fmt.Sprintf("epg-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s", tenant, application_profile)
	fvAEPg := models.NewApplicationEPG(rn, parentDn, description, fvAEPgattr)

	fvAEPg.Status = "modified"
	err := sm.Save(fvAEPg)
	return fvAEPg, err

}

func (sm *ServiceManager) ListApplicationEPG(application_profile string, tenant string) ([]*models.ApplicationEPG, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/fvAEPg.json", baseurlStr, tenant, application_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ApplicationEPGListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) SetupCreateRelationfvRsBdFromApplicationEPG(parentDn, tnFvBDName string) []byte {
	dn := fmt.Sprintf("%s/rsbd", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvBDName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsBd", dn, tnFvBDName))

	return containerJSON
}

func (sm *ServiceManager) RenderRelationfvRsAllFromApplicationEPG(fvAEPg *models.ApplicationEPG, fvRsAllData [][]byte) error {
	dn := fmt.Sprintf("%s", fvAEPg.DistinguishedName)
	headerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"		
			},
			"children":
			[`, "fvAEPg", dn))

	containerJSON := []byte{}

	fvRsItemIdx := 0
	fvRsAllDataItems := len(fvRsAllData)

	for {
		containerJSON = append(containerJSON, headerJSON...)
		sizeOfPOST := len(containerJSON) + 3

		for next := true; next; next = fvRsItemIdx < fvRsAllDataItems {
			if sizeOfPOST+len(fvRsAllData[fvRsItemIdx]) <= maxSizeOfPOST {
				containerJSON = append(containerJSON, fvRsAllData[fvRsItemIdx]...)
				containerJSON = append(containerJSON, ',')
				sizeOfPOST = len(containerJSON)
				fvRsItemIdx++
			} else {
				break
			}
		}

		containerJSON[sizeOfPOST-1] = ']'
		containerJSON = append(containerJSON, "}}"...)

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

		if fvRsItemIdx == fvRsAllDataItems {
			break
		}

		containerJSON = nil
	}

	return nil
}

func (sm *ServiceManager) CreateRelationfvRsBdFromApplicationEPG(parentDn, tnFvBDName string) error {
	dn := fmt.Sprintf("%s/rsbd", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvBDName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsBd", dn, tnFvBDName))

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

func (sm *ServiceManager) ReadRelationfvRsAllEPG(parentDn string) (map[string]interface{}, error) {
	baseurlStr := "/api/node/class"
	//dnUrl := fmt.Sprintf("%s/%s.json?rsp-subtree=children&query-target-filter=eq(fvAEPg.dn,\"%s\")", baseurlStr, "fvAEPg", parentDn)
	dnUrl := fmt.Sprintf("%s/%s.json?rsp-subtree=children&query-target-filter=eq(fvAEPg.dn,\"%s\")%s", baseurlStr, "fvAEPg", parentDn, fvRsClassesEPGFilter)
	cont, err := sm.GetViaURL(dnUrl)
	if err != nil {
		return nil, err
	}
	contList := models.ListFromContainer2(cont, "fvAEPg")

	fvRsChildren := make(map[string]interface{})

	for _, fvRsClass := range fvRsClassesEPG {
		fvRsBlock := contList[0].S(fvRsClass.Id, "attributes")
		if fvRsBlock != nil {
			if fvRsClass.TypeSet {
				fvRsBlockLen := len((fvRsBlock.Data()).([]interface{}))
				st := &schema.Set{
					F: schema.HashString,
				}
				for i := 0; i < fvRsBlockLen; i++ {
					dat := models.G((fvRsBlock).Index(i), "tDn")
					st.Add(dat)
				}
				fvRsChildren[fvRsClass.Id] = st
			} else {
				dat := models.G((fvRsBlock).Index(0), "tDn")
				fvRsChildren[fvRsClass.Id] = dat
			}
		}
	}

	return fvRsChildren, err
}

func (sm *ServiceManager) ReadRelationfvRsBdFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBd")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBd")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

func (sm *ServiceManager) SetupCreateRelationfvRsCustQosPolFromApplicationEPG(parentDn, tnQosCustomPolName string) []byte {
	dn := fmt.Sprintf("%s/rscustQosPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosCustomPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsCustQosPol", dn, tnQosCustomPolName))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsCustQosPolFromApplicationEPG(parentDn, tnQosCustomPolName string) error {
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

func (sm *ServiceManager) ReadRelationfvRsCustQosPolFromApplicationEPG(parentDn string) (interface{}, error) {
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

func (sm *ServiceManager) SetupCreateRelationfvRsDomAttFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rsdomAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsDomAtt", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsDomAttFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsdomAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsDomAtt", dn))

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

func (sm *ServiceManager) SetupDeleteRelationfvRsDomAttFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rsdomAtt-[%s]", parentDn, tDn)
	return sm.SetupDeleteByDn(dn, "fvRsDomAtt")
}

func (sm *ServiceManager) DeleteRelationfvRsDomAttFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsdomAtt-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsDomAtt")
}

func (sm *ServiceManager) ReadRelationfvRsDomAttFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsDomAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsDomAtt")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

func (sm *ServiceManager) SetupCreateRelationfvRsFcPathAttFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rsfcPathAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"		
			}
		}
	}`, "fvRsFcPathAtt", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsFcPathAttFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsfcPathAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"		
			}
		}
	}`, "fvRsFcPathAtt", dn))

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

func (sm *ServiceManager) SetupDeleteRelationfvRsFcPathAttFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rsfcPathAtt-[%s]", parentDn, tDn)
	return sm.SetupDeleteByDn(dn, "fvRsFcPathAtt")
}

func (sm *ServiceManager) DeleteRelationfvRsFcPathAttFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsfcPathAtt-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsFcPathAtt")
}

func (sm *ServiceManager) ReadRelationfvRsFcPathAttFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsFcPathAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsFcPathAtt")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

func (sm *ServiceManager) SetupCreateRelationfvRsProvFromApplicationEPG(parentDn, tnVzBrCPName string) []byte {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsProv", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsProvFromApplicationEPG(parentDn, tnVzBrCPName string) error {
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

func (sm *ServiceManager) SetupDeleteRelationfvRsProvFromApplicationEPG(parentDn, tnVzBrCPName string) []byte {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	return sm.SetupDeleteByDn(dn, "fvRsProv")
}

func (sm *ServiceManager) DeleteRelationfvRsProvFromApplicationEPG(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsProv")
}

func (sm *ServiceManager) ReadRelationfvRsProvFromApplicationEPG(parentDn string) (interface{}, error) {
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

func (sm *ServiceManager) SetupCreateRelationfvRsGraphDefFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rsgraphDef-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"				
			}
		}
	}`, "fvRsGraphDef", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsGraphDefFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsgraphDef-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"				
			}
		}
	}`, "fvRsGraphDef", dn))

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

func (sm *ServiceManager) ReadRelationfvRsGraphDefFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsGraphDef")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsGraphDef")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

func (sm *ServiceManager) SetupCreateRelationfvRsConsIfFromApplicationEPG(parentDn, tnVzCPIfName string) []byte {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsConsIf", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsConsIfFromApplicationEPG(parentDn, tnVzCPIfName string) error {
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

func (sm *ServiceManager) SetupDeleteRelationfvRsConsIfFromApplicationEPG(parentDn, tnVzCPIfName string) []byte {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	return sm.SetupDeleteByDn(dn, "fvRsConsIf")
}

func (sm *ServiceManager) DeleteRelationfvRsConsIfFromApplicationEPG(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	return sm.DeleteByDn(dn, "fvRsConsIf")
}

func (sm *ServiceManager) ReadRelationfvRsConsIfFromApplicationEPG(parentDn string) (interface{}, error) {
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

func (sm *ServiceManager) SetupCreateRelationfvRsSecInheritedFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsSecInherited", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsSecInheritedFromApplicationEPG(parentDn, tDn string) error {
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

func (sm *ServiceManager) SetupDeleteRelationfvRsSecInheritedFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	return sm.SetupDeleteByDn(dn, "fvRsSecInherited")
}

func (sm *ServiceManager) DeleteRelationfvRsSecInheritedFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsSecInherited")
}

func (sm *ServiceManager) ReadRelationfvRsSecInheritedFromApplicationEPG(parentDn string) (interface{}, error) {
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

func (sm *ServiceManager) SetupCreateRelationfvRsNodeAttFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rsnodeAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsNodeAtt", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsNodeAttFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsnodeAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsNodeAtt", dn))

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

func (sm *ServiceManager) SetupDeleteRelationfvRsNodeAttFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rsnodeAtt-[%s]", parentDn, tDn)
	return sm.SetupDeleteByDn(dn, "fvRsNodeAtt")
}

func (sm *ServiceManager) DeleteRelationfvRsNodeAttFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsnodeAtt-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsNodeAtt")
}

func (sm *ServiceManager) ReadRelationfvRsNodeAttFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsNodeAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsNodeAtt")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

func (sm *ServiceManager) SetupCreateRelationfvRsDppPolFromApplicationEPG(parentDn, tnQosDppPolName string) []byte {
	dn := fmt.Sprintf("%s/rsdppPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosDppPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsDppPol", dn, tnQosDppPolName))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsDppPolFromApplicationEPG(parentDn, tnQosDppPolName string) error {
	dn := fmt.Sprintf("%s/rsdppPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnQosDppPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsDppPol", dn, tnQosDppPolName))

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

func (sm *ServiceManager) SetupDeleteRelationfvRsDppPolFromApplicationEPG(parentDn string) []byte {
	dn := fmt.Sprintf("%s/rsdppPol", parentDn)
	return sm.SetupDeleteByDn(dn, "fvRsDppPol")
}

func (sm *ServiceManager) DeleteRelationfvRsDppPolFromApplicationEPG(parentDn string) error {
	dn := fmt.Sprintf("%s/rsdppPol", parentDn)
	return sm.DeleteByDn(dn, "fvRsDppPol")
}

func (sm *ServiceManager) ReadRelationfvRsDppPolFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsDppPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsDppPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

func (sm *ServiceManager) SetupCreateRelationfvRsConsFromApplicationEPG(parentDn, tnVzBrCPName string) []byte {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsCons", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsConsFromApplicationEPG(parentDn, tnVzBrCPName string) error {
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

func (sm *ServiceManager) SetupDeleteRelationfvRsConsFromApplicationEPG(parentDn, tnVzBrCPName string) []byte {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	return sm.SetupDeleteByDn(dn, "fvRsCons")
}

func (sm *ServiceManager) DeleteRelationfvRsConsFromApplicationEPG(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsCons")
}

func (sm *ServiceManager) ReadRelationfvRsConsFromApplicationEPG(parentDn string) (interface{}, error) {
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

func (sm *ServiceManager) SetupCreateRelationfvRsProvDefFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rsprovDef-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"				
			}
		}
	}`, "fvRsProvDef", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsProvDefFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsprovDef-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"				
			}
		}
	}`, "fvRsProvDef", dn))

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

func (sm *ServiceManager) ReadRelationfvRsProvDefFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsProvDef")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsProvDef")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

func (sm *ServiceManager) SetupCreateRelationfvRsTrustCtrlFromApplicationEPG(parentDn, tnFhsTrustCtrlPolName string) []byte {
	dn := fmt.Sprintf("%s/rstrustCtrl", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFhsTrustCtrlPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsTrustCtrl", dn, tnFhsTrustCtrlPolName))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsTrustCtrlFromApplicationEPG(parentDn, tnFhsTrustCtrlPolName string) error {
	dn := fmt.Sprintf("%s/rstrustCtrl", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFhsTrustCtrlPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsTrustCtrl", dn, tnFhsTrustCtrlPolName))

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

func (sm *ServiceManager) SetupDeleteRelationfvRsTrustCtrlFromApplicationEPG(parentDn string) []byte {
	dn := fmt.Sprintf("%s/rstrustCtrl", parentDn)
	return sm.SetupDeleteByDn(dn, "fvRsTrustCtrl")
}

func (sm *ServiceManager) DeleteRelationfvRsTrustCtrlFromApplicationEPG(parentDn string) error {
	dn := fmt.Sprintf("%s/rstrustCtrl", parentDn)
	return sm.DeleteByDn(dn, "fvRsTrustCtrl")
}

func (sm *ServiceManager) ReadRelationfvRsTrustCtrlFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsTrustCtrl")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsTrustCtrl")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

func (sm *ServiceManager) SetupCreateRelationfvRsPathAttFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rspathAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsPathAtt", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsPathAttFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rspathAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsPathAtt", dn))

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

func (sm *ServiceManager) SetupDeleteRelationfvRsPathAttFromApplicationEPG(parentDn, tDn string) []byte {
	dn := fmt.Sprintf("%s/rspathAtt-[%s]", parentDn, tDn)
	return sm.SetupDeleteByDn(dn, "fvRsPathAtt")
}

func (sm *ServiceManager) DeleteRelationfvRsPathAttFromApplicationEPG(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rspathAtt-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsPathAtt")
}

func (sm *ServiceManager) ReadRelationfvRsPathAttFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsPathAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsPathAtt")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

func (sm *ServiceManager) SetupCreateRelationfvRsProtByFromApplicationEPG(parentDn, tnVzTabooName string) []byte {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsProtBy", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsProtByFromApplicationEPG(parentDn, tnVzTabooName string) error {
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

func (sm *ServiceManager) SetupDeleteRelationfvRsProtByFromApplicationEPG(parentDn, tnVzTabooName string) []byte {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	return sm.SetupDeleteByDn(dn, "fvRsProtBy")
}

func (sm *ServiceManager) DeleteRelationfvRsProtByFromApplicationEPG(parentDn, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	return sm.DeleteByDn(dn, "fvRsProtBy")
}

func (sm *ServiceManager) ReadRelationfvRsProtByFromApplicationEPG(parentDn string) (interface{}, error) {
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

func (sm *ServiceManager) SetupCreateRelationfvRsAEPgMonPolFromApplicationEPG(parentDn, tnMonEPGPolName string) []byte {
	dn := fmt.Sprintf("%s/rsAEPgMonPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMonEPGPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsAEPgMonPol", dn, tnMonEPGPolName))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsAEPgMonPolFromApplicationEPG(parentDn, tnMonEPGPolName string) error {
	dn := fmt.Sprintf("%s/rsAEPgMonPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMonEPGPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsAEPgMonPol", dn, tnMonEPGPolName))

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

func (sm *ServiceManager) SetupDeleteRelationfvRsAEPgMonPolFromApplicationEPG(parentDn string) []byte {
	dn := fmt.Sprintf("%s/rsAEPgMonPol", parentDn)
	return sm.SetupDeleteByDn(dn, "fvRsAEPgMonPol")
}

func (sm *ServiceManager) DeleteRelationfvRsAEPgMonPolFromApplicationEPG(parentDn string) error {
	dn := fmt.Sprintf("%s/rsAEPgMonPol", parentDn)
	return sm.DeleteByDn(dn, "fvRsAEPgMonPol")
}

func (sm *ServiceManager) ReadRelationfvRsAEPgMonPolFromApplicationEPG(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsAEPgMonPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsAEPgMonPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

func (sm *ServiceManager) SetupCreateRelationfvRsIntraEpgFromApplicationEPG(parentDn, tnVzBrCPName string) []byte {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsIntraEpg", dn))

	return containerJSON
}

func (sm *ServiceManager) CreateRelationfvRsIntraEpgFromApplicationEPG(parentDn, tnVzBrCPName string) error {
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

func (sm *ServiceManager) SetupDeleteRelationfvRsIntraEpgFromApplicationEPG(parentDn, tnVzBrCPName string) []byte {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	return sm.SetupDeleteByDn(dn, "fvRsIntraEpg")
}

func (sm *ServiceManager) DeleteRelationfvRsIntraEpgFromApplicationEPG(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsIntraEpg")
}

func (sm *ServiceManager) ReadRelationfvRsIntraEpgFromApplicationEPG(parentDn string) (interface{}, error) {
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
