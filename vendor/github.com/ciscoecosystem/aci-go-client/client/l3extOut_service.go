package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateL3Outside(name string, tenant string, description string, l3extOutattr models.L3OutsideAttributes) (*models.L3Outside, error) {
	rn := fmt.Sprintf("out-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	l3extOut := models.NewL3Outside(rn, parentDn, description, l3extOutattr)
	err := sm.Save(l3extOut)
	return l3extOut, err
}

func (sm *ServiceManager) ReadL3Outside(name string, tenant string) (*models.L3Outside, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extOut := models.L3OutsideFromContainer(cont)
	return l3extOut, nil
}

func (sm *ServiceManager) DeleteL3Outside(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s", tenant, name)
	return sm.DeleteByDn(dn, models.L3extoutClassName)
}

func (sm *ServiceManager) UpdateL3Outside(name string, tenant string, description string, l3extOutattr models.L3OutsideAttributes) (*models.L3Outside, error) {
	rn := fmt.Sprintf("out-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	l3extOut := models.NewL3Outside(rn, parentDn, description, l3extOutattr)

	l3extOut.Status = "modified"
	err := sm.Save(l3extOut)
	return l3extOut, err

}

func (sm *ServiceManager) ListL3Outside(tenant string) ([]*models.L3Outside, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/l3extOut.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3OutsideListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationl3extRsDampeningPolFromL3Outside(parentDn, tnRtctrlProfileName, af string) error {
	dn := fmt.Sprintf("%s/rsdampeningPol-[%s]-%s", parentDn, tnRtctrlProfileName, af)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "l3extRsDampeningPol", dn))

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

func (sm *ServiceManager) DeleteRelationl3extRsDampeningPolFromL3Outside(parentDn, tnRtctrlProfileName, af string) error {
	dn := fmt.Sprintf("%s/rsdampeningPol-[%s]-%s", parentDn, tnRtctrlProfileName, af)
	return sm.DeleteByDn(dn, "l3extRsDampeningPol")
}

func (sm *ServiceManager) ReadRelationl3extRsDampeningPolFromL3Outside(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsDampeningPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsDampeningPol")

	st := make([]map[string]string, 0)

	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnRtctrlProfileName"] = models.G(contItem, "tDn")
		paramMap["af"] = models.G(contItem, "af")

		st = append(st, paramMap)

	}

	return st, err

}
func (sm *ServiceManager) CreateRelationl3extRsEctxFromL3Outside(parentDn, tnFvCtxName string) error {
	dn := fmt.Sprintf("%s/rsectx", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvCtxName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsEctx", dn, tnFvCtxName))

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

func (sm *ServiceManager) ReadRelationl3extRsEctxFromL3Outside(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsEctx")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsEctx")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationl3extRsOutToBDPublicSubnetHolderFromL3Outside(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsoutToBDPublicSubnetHolder-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"				
			}
		}
	}`, "l3extRsOutToBDPublicSubnetHolder", dn))

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

func (sm *ServiceManager) ReadRelationl3extRsOutToBDPublicSubnetHolderFromL3Outside(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsOutToBDPublicSubnetHolder")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsOutToBDPublicSubnetHolder")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationl3extRsInterleakPolFromL3Outside(parentDn, tnRtctrlProfileName string) error {
	dn := fmt.Sprintf("%s/rsinterleakPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnRtctrlProfileName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsInterleakPol", dn, tnRtctrlProfileName))

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

func (sm *ServiceManager) DeleteRelationl3extRsInterleakPolFromL3Outside(parentDn string) error {
	dn := fmt.Sprintf("%s/rsinterleakPol", parentDn)
	return sm.DeleteByDn(dn, "l3extRsInterleakPol")
}

func (sm *ServiceManager) ReadRelationl3extRsInterleakPolFromL3Outside(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsInterleakPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsInterleakPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationl3extRsL3DomAttFromL3Outside(parentDn, tnExtnwDomPName string) error {
	dn := fmt.Sprintf("%s/rsl3DomAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l3extRsL3DomAtt", dn, tnExtnwDomPName))

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

func (sm *ServiceManager) DeleteRelationl3extRsL3DomAttFromL3Outside(parentDn string) error {
	dn := fmt.Sprintf("%s/rsl3DomAtt", parentDn)
	return sm.DeleteByDn(dn, "l3extRsL3DomAtt")
}

func (sm *ServiceManager) ReadRelationl3extRsL3DomAttFromL3Outside(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsL3DomAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsL3DomAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
