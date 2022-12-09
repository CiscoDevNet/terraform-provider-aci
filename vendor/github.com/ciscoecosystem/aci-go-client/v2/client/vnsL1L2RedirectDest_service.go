package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateL1L2RedirectDestTraffic(destName string, parentDn string, description string, nameAlias string, vnsL1L2RedirectDestAttr models.L1L2RedirectDestTrafficAttributes) (*models.L1L2RedirectDestTraffic, error) {
	rn := fmt.Sprintf(models.RnvnsL1L2RedirectDest, destName)
	vnsL1L2RedirectDest := models.NewL1L2RedirectDestTraffic(rn, parentDn, description, nameAlias, vnsL1L2RedirectDestAttr)
	err := sm.Save(vnsL1L2RedirectDest)
	return vnsL1L2RedirectDest, err
}

func (sm *ServiceManager) ReadL1L2RedirectDestTraffic(destName string, parentDn string) (*models.L1L2RedirectDestTraffic, error) {
	dn := fmt.Sprintf("%s/%s", parentDn, fmt.Sprintf(models.RnvnsL1L2RedirectDest, destName))
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsL1L2RedirectDest := models.L1L2RedirectDestTrafficFromContainer(cont)
	return vnsL1L2RedirectDest, nil
}

func (sm *ServiceManager) DeleteL1L2RedirectDestTraffic(destName string, parentDn string) error {
	dn := fmt.Sprintf("%s/%s", parentDn, fmt.Sprintf(models.RnvnsL1L2RedirectDest, destName))
	return sm.DeleteByDn(dn, models.Vnsl1l2redirectdestClassName)
}

func (sm *ServiceManager) UpdateL1L2RedirectDestTraffic(destName string, parentDn string, description string, nameAlias string, vnsL1L2RedirectDestAttr models.L1L2RedirectDestTrafficAttributes) (*models.L1L2RedirectDestTraffic, error) {
	rn := fmt.Sprintf(models.RnvnsL1L2RedirectDest, destName)
	vnsL1L2RedirectDest := models.NewL1L2RedirectDestTraffic(rn, parentDn, description, nameAlias, vnsL1L2RedirectDestAttr)
	vnsL1L2RedirectDest.Status = "modified"
	err := sm.Save(vnsL1L2RedirectDest)
	return vnsL1L2RedirectDest, err
}

func (sm *ServiceManager) ListL1L2RedirectDestTraffic(parentDn string) ([]*models.L1L2RedirectDestTraffic, error) {
	dnUrl := fmt.Sprintf("%s/%s/vnsL1L2RedirectDest.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.L1L2RedirectDestTrafficListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsL1L2RedirectHealthGroup(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsL1L2RedirectHealthGroup", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsL1L2RedirectHealthGroup", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvnsRsL1L2RedirectHealthGroup(parentDn string) error {
	dn := fmt.Sprintf("%s/rsL1L2RedirectHealthGroup", parentDn)
	return sm.DeleteByDn(dn, "vnsRsL1L2RedirectHealthGroup")
}

func (sm *ServiceManager) ReadRelationvnsRsL1L2RedirectHealthGroup(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsL1L2RedirectHealthGroup")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsL1L2RedirectHealthGroup")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvnsRsToCIf(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rstoCIf", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsToCIf", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvnsRsToCIf(parentDn string) error {
	dn := fmt.Sprintf("%s/rstoCIf", parentDn)
	return sm.DeleteByDn(dn, "vnsRsToCIf")
}

func (sm *ServiceManager) ReadRelationvnsRsToCIf(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsToCIf")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsToCIf")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
