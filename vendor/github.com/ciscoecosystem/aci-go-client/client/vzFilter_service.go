package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFilter(name string, tenant string, description string, vzFilterattr models.FilterAttributes) (*models.Filter, error) {
	rn := fmt.Sprintf("flt-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vzFilter := models.NewFilter(rn, parentDn, description, vzFilterattr)
	err := sm.Save(vzFilter)
	return vzFilter, err
}

func (sm *ServiceManager) ReadFilter(name string, tenant string) (*models.Filter, error) {
	dn := fmt.Sprintf("uni/tn-%s/flt-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzFilter := models.FilterFromContainer(cont)
	return vzFilter, nil
}

func (sm *ServiceManager) DeleteFilter(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/flt-%s", tenant, name)
	return sm.DeleteByDn(dn, models.VzfilterClassName)
}

func (sm *ServiceManager) UpdateFilter(name string, tenant string, description string, vzFilterattr models.FilterAttributes) (*models.Filter, error) {
	rn := fmt.Sprintf("flt-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vzFilter := models.NewFilter(rn, parentDn, description, vzFilterattr)

	vzFilter.Status = "modified"
	err := sm.Save(vzFilter)
	return vzFilter, err

}

func (sm *ServiceManager) ListFilter(tenant string) ([]*models.Filter, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/vzFilter.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FilterListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvzRsFiltGraphAttFromFilter(parentDn, tnVnsInTermName string) error {
	dn := fmt.Sprintf("%s/rsFiltGraphAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVnsInTermName": "%s"
								
			}
		}
	}`, "vzRsFiltGraphAtt", dn, tnVnsInTermName))

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

func (sm *ServiceManager) ReadRelationvzRsFiltGraphAttFromFilter(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsFiltGraphAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsFiltGraphAtt")

	if len(contList) > 0 {
		dat := models.CurlyBraces(models.G(contList[0], "tnVnsInTermName"))
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvzRsFwdRFltPAttFromFilter(parentDn, tnVzAFilterableUnitName string) error {
	dn := fmt.Sprintf("%s/rsFwdRFltPAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVzAFilterableUnitName": "%s"
								
			}
		}
	}`, "vzRsFwdRFltPAtt", dn, tnVzAFilterableUnitName))

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

func (sm *ServiceManager) ReadRelationvzRsFwdRFltPAttFromFilter(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsFwdRFltPAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsFwdRFltPAtt")

	if len(contList) > 0 {
		dat := models.CurlyBraces(models.G(contList[0], "tnVzAFilterableUnitName"))
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvzRsRevRFltPAttFromFilter(parentDn, tnVzAFilterableUnitName string) error {
	dn := fmt.Sprintf("%s/rsRevRFltPAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVzAFilterableUnitName": "%s"
								
			}
		}
	}`, "vzRsRevRFltPAtt", dn, tnVzAFilterableUnitName))

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

func (sm *ServiceManager) ReadRelationvzRsRevRFltPAttFromFilter(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsRevRFltPAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsRevRFltPAtt")

	if len(contList) > 0 {
		dat := models.CurlyBraces(models.G(contList[0], "tnVzAFilterableUnitName"))
		return dat, err
	} else {
		return nil, err
	}

}
