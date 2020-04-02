package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFCDomain(name string, description string, fcDomPattr models.FCDomainAttributes) (*models.FCDomain, error) {
	rn := fmt.Sprintf("fc-%s", name)
	parentDn := fmt.Sprintf("uni")
	fcDomP := models.NewFCDomain(rn, parentDn, description, fcDomPattr)
	err := sm.Save(fcDomP)
	return fcDomP, err
}

func (sm *ServiceManager) ReadFCDomain(name string) (*models.FCDomain, error) {
	dn := fmt.Sprintf("uni/fc-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fcDomP := models.FCDomainFromContainer(cont)
	return fcDomP, nil
}

func (sm *ServiceManager) DeleteFCDomain(name string) error {
	dn := fmt.Sprintf("uni/fc-%s", name)
	return sm.DeleteByDn(dn, models.FcdompClassName)
}

func (sm *ServiceManager) UpdateFCDomain(name string, description string, fcDomPattr models.FCDomainAttributes) (*models.FCDomain, error) {
	rn := fmt.Sprintf("fc-%s", name)
	parentDn := fmt.Sprintf("uni")
	fcDomP := models.NewFCDomain(rn, parentDn, description, fcDomPattr)

	fcDomP.Status = "modified"
	err := sm.Save(fcDomP)
	return fcDomP, err

}

func (sm *ServiceManager) ListFCDomain() ([]*models.FCDomain, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fcDomP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FCDomainListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsVlanNsFromFCDomain(parentDn, tnFvnsVlanInstPName string) error {
	dn := fmt.Sprintf("%s/rsvlanNs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvnsVlanInstPName": "%s"
								
			}
		}
	}`, "infraRsVlanNs", dn, tnFvnsVlanInstPName))

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

func (sm *ServiceManager) DeleteRelationinfraRsVlanNsFromFCDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvlanNs", parentDn)
	return sm.DeleteByDn(dn, "infraRsVlanNs")
}

func (sm *ServiceManager) ReadRelationinfraRsVlanNsFromFCDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsVlanNs")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsVlanNs")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFvnsVlanInstPName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfcRsVsanNsFromFCDomain(parentDn, tnFvnsVsanInstPName string) error {
	dn := fmt.Sprintf("%s/rsvsanNs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvnsVsanInstPName": "%s"
								
			}
		}
	}`, "fcRsVsanNs", dn, tnFvnsVsanInstPName))

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

func (sm *ServiceManager) DeleteRelationfcRsVsanNsFromFCDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvsanNs", parentDn)
	return sm.DeleteByDn(dn, "fcRsVsanNs")
}

func (sm *ServiceManager) ReadRelationfcRsVsanNsFromFCDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fcRsVsanNs")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fcRsVsanNs")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFvnsVsanInstPName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfcRsVsanAttrFromFCDomain(parentDn, tnFcVsanAttrPName string) error {
	dn := fmt.Sprintf("%s/rsvsanAttr", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFcVsanAttrPName": "%s"
								
			}
		}
	}`, "fcRsVsanAttr", dn, tnFcVsanAttrPName))

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

func (sm *ServiceManager) DeleteRelationfcRsVsanAttrFromFCDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvsanAttr", parentDn)
	return sm.DeleteByDn(dn, "fcRsVsanAttr")
}

func (sm *ServiceManager) ReadRelationfcRsVsanAttrFromFCDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fcRsVsanAttr")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fcRsVsanAttr")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFcVsanAttrPName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsVlanNsDefFromFCDomain(parentDn, tnFvnsAInstPName string) error {
	dn := fmt.Sprintf("%s/rsvlanNsDef", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvnsAInstPName": "%s"
								
			}
		}
	}`, "infraRsVlanNsDef", dn, tnFvnsAInstPName))

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

func (sm *ServiceManager) ReadRelationinfraRsVlanNsDefFromFCDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsVlanNsDef")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsVlanNsDef")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFvnsAInstPName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsVipAddrNsFromFCDomain(parentDn, tnFvnsAddrInstName string) error {
	dn := fmt.Sprintf("%s/rsvipAddrNs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvnsAddrInstName": "%s"
								
			}
		}
	}`, "infraRsVipAddrNs", dn, tnFvnsAddrInstName))

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

func (sm *ServiceManager) DeleteRelationinfraRsVipAddrNsFromFCDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvipAddrNs", parentDn)
	return sm.DeleteByDn(dn, "infraRsVipAddrNs")
}

func (sm *ServiceManager) ReadRelationinfraRsVipAddrNsFromFCDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsVipAddrNs")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsVipAddrNs")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFvnsAddrInstName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsDomVxlanNsDefFromFCDomain(parentDn, tnFvnsAInstPName string) error {
	dn := fmt.Sprintf("%s/rsdomVxlanNsDef", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvnsAInstPName": "%s"
								
			}
		}
	}`, "infraRsDomVxlanNsDef", dn, tnFvnsAInstPName))

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

func (sm *ServiceManager) ReadRelationinfraRsDomVxlanNsDefFromFCDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsDomVxlanNsDef")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsDomVxlanNsDef")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFvnsAInstPName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfcRsVsanAttrDefFromFCDomain(parentDn, tnFcVsanAttrPName string) error {
	dn := fmt.Sprintf("%s/rsvsanAttrDef", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFcVsanAttrPName": "%s"
								
			}
		}
	}`, "fcRsVsanAttrDef", dn, tnFcVsanAttrPName))

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

func (sm *ServiceManager) ReadRelationfcRsVsanAttrDefFromFCDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fcRsVsanAttrDef")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fcRsVsanAttrDef")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFcVsanAttrPName")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfcRsVsanNsDefFromFCDomain(parentDn, tnFvnsAVsanInstPName string) error {
	dn := fmt.Sprintf("%s/rsvsanNsDef", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvnsAVsanInstPName": "%s"
								
			}
		}
	}`, "fcRsVsanNsDef", dn, tnFvnsAVsanInstPName))

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

func (sm *ServiceManager) ReadRelationfcRsVsanNsDefFromFCDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fcRsVsanNsDef")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fcRsVsanNsDef")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFvnsAVsanInstPName")
		return dat, err
	} else {
		return nil, err
	}

}
