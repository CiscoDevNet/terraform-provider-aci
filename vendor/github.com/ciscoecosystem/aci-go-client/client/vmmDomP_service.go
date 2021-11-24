package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateVMMDomain(name string, provider_profile_vendor string, vmmDomPattr models.VMMDomainAttributes) (*models.VMMDomain, error) {
	rn := fmt.Sprintf("dom-%s", name)
	parentDn := fmt.Sprintf("uni/vmmp-%s", provider_profile_vendor)
	vmmDomP := models.NewVMMDomain(rn, parentDn, vmmDomPattr)
	err := sm.Save(vmmDomP)
	return vmmDomP, err
}

func (sm *ServiceManager) ReadVMMDomain(name string, provider_profile_vendor string) (*models.VMMDomain, error) {
	dn := fmt.Sprintf("uni/vmmp-%s/dom-%s", provider_profile_vendor, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vmmDomP := models.VMMDomainFromContainer(cont)
	return vmmDomP, nil
}

func (sm *ServiceManager) DeleteVMMDomain(name string, provider_profile_vendor string) error {
	dn := fmt.Sprintf("uni/vmmp-%s/dom-%s", provider_profile_vendor, name)
	return sm.DeleteByDn(dn, models.VmmdompClassName)
}

func (sm *ServiceManager) UpdateVMMDomain(name string, provider_profile_vendor string, vmmDomPattr models.VMMDomainAttributes) (*models.VMMDomain, error) {
	rn := fmt.Sprintf("dom-%s", name)
	parentDn := fmt.Sprintf("uni/vmmp-%s", provider_profile_vendor)
	vmmDomP := models.NewVMMDomain(rn, parentDn, vmmDomPattr)

	vmmDomP.Status = "modified"
	err := sm.Save(vmmDomP)
	return vmmDomP, err

}

func (sm *ServiceManager) ListVMMDomain(provider_profile_vendor string) ([]*models.VMMDomain, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/vmmp-%s/vmmDomP.json", baseurlStr, provider_profile_vendor)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.VMMDomainListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvmmRsPrefEnhancedLagPolFromVMMDomain(parentDn, tnLacpEnhancedLagPolName string) error {
	dn := fmt.Sprintf("%s/rsprefEnhancedLagPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vmmRsPrefEnhancedLagPol", dn, tnLacpEnhancedLagPolName))

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

func (sm *ServiceManager) DeleteRelationvmmRsPrefEnhancedLagPolFromVMMDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsprefEnhancedLagPol", parentDn)
	return sm.DeleteByDn(dn, "vmmRsPrefEnhancedLagPol")
}

func (sm *ServiceManager) ReadRelationvmmRsPrefEnhancedLagPolFromVMMDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vmmRsPrefEnhancedLagPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vmmRsPrefEnhancedLagPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsVlanNsFromVMMDomain(parentDn, tnFvnsVlanInstPName string) error {
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

func (sm *ServiceManager) DeleteRelationinfraRsVlanNsFromVMMDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvlanNs", parentDn)
	return sm.DeleteByDn(dn, "infraRsVlanNs")
}

func (sm *ServiceManager) ReadRelationinfraRsVlanNsFromVMMDomain(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationvmmRsDomMcastAddrNsFromVMMDomain(parentDn, tnFvnsMcastAddrInstPName string) error {
	dn := fmt.Sprintf("%s/rsdomMcastAddrNs", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vmmRsDomMcastAddrNs", dn, tnFvnsMcastAddrInstPName))

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

func (sm *ServiceManager) DeleteRelationvmmRsDomMcastAddrNsFromVMMDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsdomMcastAddrNs", parentDn)
	return sm.DeleteByDn(dn, "vmmRsDomMcastAddrNs")
}

func (sm *ServiceManager) ReadRelationvmmRsDomMcastAddrNsFromVMMDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vmmRsDomMcastAddrNs")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vmmRsDomMcastAddrNs")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvmmRsDefaultCdpIfPolFromVMMDomain(parentDn, tnCdpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsdefaultCdpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnCdpIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vmmRsDefaultCdpIfPol", dn, tnCdpIfPolName))

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

func (sm *ServiceManager) ReadRelationvmmRsDefaultCdpIfPolFromVMMDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vmmRsDefaultCdpIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vmmRsDefaultCdpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvmmRsDefaultLacpLagPolFromVMMDomain(parentDn, tnLacpLagPolName string) error {
	dn := fmt.Sprintf("%s/rsdefaultLacpLagPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnLacpLagPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vmmRsDefaultLacpLagPol", dn, tnLacpLagPolName))

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

func (sm *ServiceManager) ReadRelationvmmRsDefaultLacpLagPolFromVMMDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vmmRsDefaultLacpLagPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vmmRsDefaultLacpLagPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsVlanNsDefFromVMMDomain(parentDn, tnFvnsAInstPName string) error {
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

func (sm *ServiceManager) ReadRelationinfraRsVlanNsDefFromVMMDomain(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationinfraRsVipAddrNsFromVMMDomain(parentDn, tnFvnsAddrInstName string) error {
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

func (sm *ServiceManager) DeleteRelationinfraRsVipAddrNsFromVMMDomain(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvipAddrNs", parentDn)
	return sm.DeleteByDn(dn, "infraRsVipAddrNs")
}

func (sm *ServiceManager) ReadRelationinfraRsVipAddrNsFromVMMDomain(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationvmmRsDefaultLldpIfPolFromVMMDomain(parentDn, tnLldpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsdefaultLldpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnLldpIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vmmRsDefaultLldpIfPol", dn, tnLldpIfPolName))

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

func (sm *ServiceManager) ReadRelationvmmRsDefaultLldpIfPolFromVMMDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vmmRsDefaultLldpIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vmmRsDefaultLldpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvmmRsDefaultStpIfPolFromVMMDomain(parentDn, tnStpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsdefaultStpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnStpIfPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vmmRsDefaultStpIfPol", dn, tnStpIfPolName))

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

func (sm *ServiceManager) ReadRelationvmmRsDefaultStpIfPolFromVMMDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vmmRsDefaultStpIfPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vmmRsDefaultStpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsDomVxlanNsDefFromVMMDomain(parentDn, tnFvnsAInstPName string) error {
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

func (sm *ServiceManager) ReadRelationinfraRsDomVxlanNsDefFromVMMDomain(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationvmmRsDefaultFwPolFromVMMDomain(parentDn, tnNwsFwPolName string) error {
	dn := fmt.Sprintf("%s/rsdefaultFwPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnNwsFwPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vmmRsDefaultFwPol", dn, tnNwsFwPolName))

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

func (sm *ServiceManager) ReadRelationvmmRsDefaultFwPolFromVMMDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vmmRsDefaultFwPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vmmRsDefaultFwPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvmmRsDefaultL2InstPolFromVMMDomain(parentDn, tnL2InstPolName string) error {
	dn := fmt.Sprintf("%s/rsdefaultL2InstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnL2InstPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vmmRsDefaultL2InstPol", dn, tnL2InstPolName))

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

func (sm *ServiceManager) ReadRelationvmmRsDefaultL2InstPolFromVMMDomain(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vmmRsDefaultL2InstPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vmmRsDefaultL2InstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
