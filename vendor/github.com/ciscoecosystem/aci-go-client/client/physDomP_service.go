package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreatePhysicalDomain(name string, description string, physDomPattr models.PhysicalDomainAttributes) (*models.PhysicalDomain, error) {
	rn := fmt.Sprintf("phys-%s", name)
	parentDn := fmt.Sprintf("uni")
	physDomP := models.NewPhysicalDomain(rn, parentDn, description, physDomPattr)
	err := sm.Save(physDomP)
	return physDomP, err
}

func (sm *ServiceManager) ReadPhysicalDomain(name string) (*models.PhysicalDomain, error) {
	dn := fmt.Sprintf("uni/phys-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	physDomP := models.PhysicalDomainFromContainer(cont)
	return physDomP, nil
}

func (sm *ServiceManager) DeletePhysicalDomain(name string) error {
	dn := fmt.Sprintf("uni/phys-%s", name)
	return sm.DeleteByDn(dn, models.PhysdompClassName)
}

func (sm *ServiceManager) UpdatePhysicalDomain(name string, description string, physDomPattr models.PhysicalDomainAttributes) (*models.PhysicalDomain, error) {
	rn := fmt.Sprintf("phys-%s", name)
	parentDn := fmt.Sprintf("uni")
	physDomP := models.NewPhysicalDomain(rn, parentDn, description, physDomPattr)

	physDomP.Status = "modified"
	err := sm.Save(physDomP)
	return physDomP, err

}

func (sm *ServiceManager) ListPhysicalDomain() ([]*models.PhysicalDomain, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/physDomP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.PhysicalDomainListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsVlanNsFromPhysicalDomain(parentDn, tnFvnsVlanInstPName string) error {
	dn := fmt.Sprintf("%s/rsvlanNs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) DeleteRelationinfraRsVlanNsFromPhysicalDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvlanNs", parentDn)
	return sm.DeleteByDn(dn, "infraRsVlanNs")
}

func (sm *ServiceManager) ReadRelationinfraRsVlanNsFromPhysicalDomain(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationinfraRsVlanNsDefFromPhysicalDomain(parentDn, tnFvnsAInstPName string) error {
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) ReadRelationinfraRsVlanNsDefFromPhysicalDomain(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationinfraRsVipAddrNsFromPhysicalDomain(parentDn, tnFvnsAddrInstName string) error {
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) DeleteRelationinfraRsVipAddrNsFromPhysicalDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvipAddrNs", parentDn)
	return sm.DeleteByDn(dn, "infraRsVipAddrNs")
}

func (sm *ServiceManager) ReadRelationinfraRsVipAddrNsFromPhysicalDomain(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationinfraRsDomVxlanNsDefFromPhysicalDomain(parentDn, tnFvnsAInstPName string) error {
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

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) ReadRelationinfraRsDomVxlanNsDefFromPhysicalDomain(parentDn string) (interface{}, error) {
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
