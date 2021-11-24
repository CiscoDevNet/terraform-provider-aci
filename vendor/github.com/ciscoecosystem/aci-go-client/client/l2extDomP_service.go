package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateL2Domain(name string, l2extDomPattr models.L2DomainAttributes) (*models.L2Domain, error) {
	rn := fmt.Sprintf("l2dom-%s", name)
	parentDn := fmt.Sprintf("uni")
	l2extDomP := models.NewL2Domain(rn, parentDn, l2extDomPattr)
	err := sm.Save(l2extDomP)
	return l2extDomP, err
}

func (sm *ServiceManager) ReadL2Domain(name string) (*models.L2Domain, error) {
	dn := fmt.Sprintf("uni/l2dom-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l2extDomP := models.L2DomainFromContainer(cont)
	return l2extDomP, nil
}

func (sm *ServiceManager) DeleteL2Domain(name string) error {
	dn := fmt.Sprintf("uni/l2dom-%s", name)
	return sm.DeleteByDn(dn, models.L2extdompClassName)
}

func (sm *ServiceManager) UpdateL2Domain(name string, l2extDomPattr models.L2DomainAttributes) (*models.L2Domain, error) {
	rn := fmt.Sprintf("l2dom-%s", name)
	parentDn := fmt.Sprintf("uni")
	l2extDomP := models.NewL2Domain(rn, parentDn, l2extDomPattr)

	l2extDomP.Status = "modified"
	err := sm.Save(l2extDomP)
	return l2extDomP, err

}

func (sm *ServiceManager) ListL2Domain() ([]*models.L2Domain, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/l2extDomP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L2DomainListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsVlanNsFromL2Domain(parentDn, tnFvnsVlanInstPName string) error {
	dn := fmt.Sprintf("%s/rsvlanNs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s",
				"annotation":"orchestrator:terraform"
								
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

func (sm *ServiceManager) DeleteRelationinfraRsVlanNsFromL2Domain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvlanNs", parentDn)
	return sm.DeleteByDn(dn, "infraRsVlanNs")
}

func (sm *ServiceManager) ReadRelationinfraRsVlanNsFromL2Domain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsVlanNs")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsVlanNs")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsVlanNsDefFromL2Domain(parentDn, tnFvnsAInstPName string) error {
	dn := fmt.Sprintf("%s/rsvlanNsDef", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s"
								
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

func (sm *ServiceManager) ReadRelationinfraRsVlanNsDefFromL2Domain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsVlanNsDef")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsVlanNsDef")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsVipAddrNsFromL2Domain(parentDn, tnFvnsAddrInstName string) error {
	dn := fmt.Sprintf("%s/rsvipAddrNs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
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

func (sm *ServiceManager) DeleteRelationinfraRsVipAddrNsFromL2Domain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvipAddrNs", parentDn)
	return sm.DeleteByDn(dn, "infraRsVipAddrNs")
}

func (sm *ServiceManager) ReadRelationinfraRsVipAddrNsFromL2Domain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsVipAddrNs")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsVipAddrNs")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationextnwRsOutFromL2Domain(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsout-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn":"%s"			
			}
		}
	}`, "extnwRsOut", dn, tDn))

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

func (sm *ServiceManager) ReadRelationextnwRsOutFromL2Domain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "extnwRsOut")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "extnwRsOut")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationinfraRsDomVxlanNsDefFromL2Domain(parentDn, tnFvnsAInstPName string) error {
	dn := fmt.Sprintf("%s/rsdomVxlanNsDef", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s"
								
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

func (sm *ServiceManager) ReadRelationinfraRsDomVxlanNsDefFromL2Domain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsDomVxlanNsDef")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsDomVxlanNsDef")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
