package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateBridgeDomain(name string, tenant string, description string, fvBDattr models.BridgeDomainAttributes) (*models.BridgeDomain, error) {
	rn := fmt.Sprintf("BD-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	fvBD := models.NewBridgeDomain(rn, parentDn, description, fvBDattr)
	err := sm.Save(fvBD)
	return fvBD, err
}

func (sm *ServiceManager) ReadBridgeDomain(name string, tenant string) (*models.BridgeDomain, error) {
	dn := fmt.Sprintf("uni/tn-%s/BD-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvBD := models.BridgeDomainFromContainer(cont)
	return fvBD, nil
}

func (sm *ServiceManager) DeleteBridgeDomain(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/BD-%s", tenant, name)
	return sm.DeleteByDn(dn, models.FvbdClassName)
}

func (sm *ServiceManager) UpdateBridgeDomain(name string, tenant string, description string, fvBDattr models.BridgeDomainAttributes) (*models.BridgeDomain, error) {
	rn := fmt.Sprintf("BD-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	fvBD := models.NewBridgeDomain(rn, parentDn, description, fvBDattr)

	fvBD.Status = "modified"
	err := sm.Save(fvBD)
	return fvBD, err

}

func (sm *ServiceManager) ListBridgeDomain(tenant string) ([]*models.BridgeDomain, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/fvBD.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BridgeDomainListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationfvRsBDToProfileFromBridgeDomain(parentDn, tnRtctrlProfileName string) error {
	dn := fmt.Sprintf("%s/rsBDToProfile", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnRtctrlProfileName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsBDToProfile", dn, tnRtctrlProfileName))

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

func (sm *ServiceManager) DeleteRelationfvRsBDToProfileFromBridgeDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsBDToProfile", parentDn)
	return sm.DeleteByDn(dn, "fvRsBDToProfile")
}

func (sm *ServiceManager) ReadRelationfvRsBDToProfileFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBDToProfile")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBDToProfile")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsMldsnFromBridgeDomain(parentDn, tnMldSnoopPolName string) error {
	dn := fmt.Sprintf("%s/rsmldsn", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMldSnoopPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsMldsn", dn, tnMldSnoopPolName))

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

func (sm *ServiceManager) ReadRelationfvRsMldsnFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsMldsn")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsMldsn")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsABDPolMonPolFromBridgeDomain(parentDn, tnMonEPGPolName string) error {
	dn := fmt.Sprintf("%s/rsABDPolMonPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMonEPGPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsABDPolMonPol", dn, tnMonEPGPolName))

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

func (sm *ServiceManager) DeleteRelationfvRsABDPolMonPolFromBridgeDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsABDPolMonPol", parentDn)
	return sm.DeleteByDn(dn, "fvRsABDPolMonPol")
}

func (sm *ServiceManager) ReadRelationfvRsABDPolMonPolFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsABDPolMonPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsABDPolMonPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsBDToNdPFromBridgeDomain(parentDn, tnNdIfPolName string) error {
	dn := fmt.Sprintf("%s/rsBDToNdP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnNdIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsBDToNdP", dn, tnNdIfPolName))

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

func (sm *ServiceManager) ReadRelationfvRsBDToNdPFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBDToNdP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBDToNdP")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsBdFloodToFromBridgeDomain(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsbdFloodTo-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsBdFloodTo", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsBdFloodToFromBridgeDomain(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsbdFloodTo-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsBdFloodTo")
}

func (sm *ServiceManager) ReadRelationfvRsBdFloodToFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBdFloodTo")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBdFloodTo")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsBDToFhsFromBridgeDomain(parentDn, tnFhsBDPolName string) error {
	dn := fmt.Sprintf("%s/rsBDToFhs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFhsBDPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsBDToFhs", dn, tnFhsBDPolName))

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

func (sm *ServiceManager) DeleteRelationfvRsBDToFhsFromBridgeDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsBDToFhs", parentDn)
	return sm.DeleteByDn(dn, "fvRsBDToFhs")
}

func (sm *ServiceManager) ReadRelationfvRsBDToFhsFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBDToFhs")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBDToFhs")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsBDToRelayPFromBridgeDomain(parentDn, tnDhcpRelayPName string) error {
	dn := fmt.Sprintf("%s/rsBDToRelayP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnDhcpRelayPName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsBDToRelayP", dn, tnDhcpRelayPName))

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

func (sm *ServiceManager) DeleteRelationfvRsBDToRelayPFromBridgeDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsBDToRelayP", parentDn)
	return sm.DeleteByDn(dn, "fvRsBDToRelayP")
}

func (sm *ServiceManager) ReadRelationfvRsBDToRelayPFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBDToRelayP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBDToRelayP")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsCtxFromBridgeDomain(parentDn, tnFvCtxName string) error {
	dn := fmt.Sprintf("%s/rsctx", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvCtxName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsCtx", dn, tnFvCtxName))

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

func (sm *ServiceManager) ReadRelationfvRsCtxFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsCtx")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsCtx")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(parentDn, tnNetflowMonitorPolName, fltType string) error {
	dn := fmt.Sprintf("%s/rsBDToNetflowMonitorPol-[%s]-%s", parentDn, tnNetflowMonitorPolName, fltType)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsBDToNetflowMonitorPol", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(parentDn, tnNetflowMonitorPolName, fltType string) error {
	dn := fmt.Sprintf("%s/rsBDToNetflowMonitorPol-[%s]-%s", parentDn, tnNetflowMonitorPolName, fltType)
	return sm.DeleteByDn(dn, "fvRsBDToNetflowMonitorPol")
}

func (sm *ServiceManager) ReadRelationfvRsBDToNetflowMonitorPolFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBDToNetflowMonitorPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBDToNetflowMonitorPol")

	st := make([]map[string]string, 0)

	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["tnNetflowMonitorPolName"] = models.G(contItem, "tDn")
		paramMap["fltType"] = models.G(contItem, "fltType")

		st = append(st, paramMap)

	}

	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsIgmpsnFromBridgeDomain(parentDn, tnIgmpSnoopPolName string) error {
	dn := fmt.Sprintf("%s/rsigmpsn", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnIgmpSnoopPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsIgmpsn", dn, tnIgmpSnoopPolName))

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

func (sm *ServiceManager) ReadRelationfvRsIgmpsnFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsIgmpsn")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsIgmpsn")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsBdToEpRetFromBridgeDomain(parentDn, tnFvEpRetPolName string) error {
	dn := fmt.Sprintf("%s/rsbdToEpRet", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvEpRetPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsBdToEpRet", dn, tnFvEpRetPolName))

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

func (sm *ServiceManager) ReadRelationfvRsBdToEpRetFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBdToEpRet")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBdToEpRet")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsBDToOutFromBridgeDomain(parentDn, tnL3extOutName string) error {
	dn := fmt.Sprintf("%s/rsBDToOut-%s", parentDn, tnL3extOutName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsBDToOut", dn))

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

func (sm *ServiceManager) DeleteRelationfvRsBDToOutFromBridgeDomain(parentDn, tnL3extOutName string) error {
	dn := fmt.Sprintf("%s/rsBDToOut-%s", parentDn, tnL3extOutName)
	return sm.DeleteByDn(dn, "fvRsBDToOut")
}

func (sm *ServiceManager) ReadRelationfvRsBDToOutFromBridgeDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBDToOut")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBDToOut")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
