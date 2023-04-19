package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateDestinationofredirectedtraffic(ip string, parentDn string, description string, nameAlias string, vnsRedirectDestAttr models.DestinationofredirectedtrafficAttributes) (*models.Destinationofredirectedtraffic, error) {
	rn := fmt.Sprintf(models.RnvnsRedirectDest, ip)
	vnsRedirectDest := models.NewDestinationofredirectedtraffic(rn, parentDn, description, nameAlias, vnsRedirectDestAttr)
	err := sm.Save(vnsRedirectDest)
	return vnsRedirectDest, err
}

func (sm *ServiceManager) ReadDestinationofredirectedtraffic(ip string, parentDn string) (*models.Destinationofredirectedtraffic, error) {
	dn := fmt.Sprintf("%s/%s", parentDn, fmt.Sprintf(models.RnvnsRedirectDest, ip))

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsRedirectDest := models.DestinationofredirectedtrafficFromContainer(cont)
	return vnsRedirectDest, nil
}

func (sm *ServiceManager) DeleteDestinationofredirectedtraffic(ip string, parentDn string) error {
	dn := fmt.Sprintf("%s/%s", parentDn, fmt.Sprintf(models.RnvnsRedirectDest, ip))
	return sm.DeleteByDn(dn, models.VnsredirectdestClassName)
}

func (sm *ServiceManager) UpdateDestinationofredirectedtraffic(ip string, parentDn string, description string, nameAlias string, vnsRedirectDestAttr models.DestinationofredirectedtrafficAttributes) (*models.Destinationofredirectedtraffic, error) {
	rn := fmt.Sprintf(models.RnvnsRedirectDest, ip)
	vnsRedirectDest := models.NewDestinationofredirectedtraffic(rn, parentDn, description, nameAlias, vnsRedirectDestAttr)
	vnsRedirectDest.Status = "modified"
	err := sm.Save(vnsRedirectDest)
	return vnsRedirectDest, err
}

func (sm *ServiceManager) ListDestinationofredirectedtraffic(parentDn string) ([]*models.Destinationofredirectedtraffic, error) {
	dnUrl := fmt.Sprintf("%s/%s/vnsRedirectDest.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.DestinationofredirectedtrafficListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsRedirectHealthGroup(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsRedirectHealthGroup", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsRedirectHealthGroup", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvnsRsRedirectHealthGroup(parentDn string) error {
	dn := fmt.Sprintf("%s/rsRedirectHealthGroup", parentDn)
	return sm.DeleteByDn(dn, "vnsRsRedirectHealthGroup")
}

func (sm *ServiceManager) ReadRelationvnsRsRedirectHealthGroup(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsRedirectHealthGroup")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsRedirectHealthGroup")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
