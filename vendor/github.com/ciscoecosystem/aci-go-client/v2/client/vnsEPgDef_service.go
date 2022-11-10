package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateEPgDef(name string, parent_dn string, description string, nameAlias string, vnsEPgDefAttr models.EPgDefAttributes) (*models.EPgDef, error) {
	rn := fmt.Sprintf(models.RnvnsEPgDef, name)
	vnsEPgDef := models.NewEPgDef(rn, parent_dn, description, nameAlias, vnsEPgDefAttr)
	err := sm.Save(vnsEPgDef)
	return vnsEPgDef, err
}

func (sm *ServiceManager) ReadEPgDef(name string, parent_dn string) (*models.EPgDef, error) {
	dn := fmt.Sprintf(models.DnvnsEPgDef, parent_dn, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsEPgDef := models.EPgDefFromContainer(cont)
	return vnsEPgDef, nil
}

func (sm *ServiceManager) DeleteEPgDef(name string, parent_dn string) error {
	dn := fmt.Sprintf(models.DnvnsEPgDef, parent_dn, name)
	return sm.DeleteByDn(dn, models.VnsepgdefClassName)
}

func (sm *ServiceManager) UpdateEPgDef(name string, parent_dn string, description string, nameAlias string, vnsEPgDefAttr models.EPgDefAttributes) (*models.EPgDef, error) {
	rn := fmt.Sprintf(models.RnvnsEPgDef, name)
	vnsEPgDef := models.NewEPgDef(rn, parent_dn, description, nameAlias, vnsEPgDefAttr)
	vnsEPgDef.Status = "modified"
	err := sm.Save(vnsEPgDef)
	return vnsEPgDef, err
}

func (sm *ServiceManager) ListEPgDef(parent_dn string) ([]*models.EPgDef, error) {
	dnUrl := fmt.Sprintf("%s/%s/vnsEPgDef.json", models.BaseurlStr, parent_dn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.EPgDefListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsEPgDefToConn(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsEPgDefToConn", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsEPgDefToConn", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvnsRsEPgDefToConn(parentDn string) error {
	dn := fmt.Sprintf("%s/rsEPgDefToConn", parentDn)
	return sm.DeleteByDn(dn, "vnsRsEPgDefToConn")
}

func (sm *ServiceManager) ReadRelationvnsRsEPgDefToConn(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsEPgDefToConn")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsEPgDefToConn")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvnsRsEPgDefToLIf(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsEPgDefToLIf", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsEPgDefToLIf", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvnsRsEPgDefToLIf(parentDn string) error {
	dn := fmt.Sprintf("%s/rsEPgDefToLIf", parentDn)
	return sm.DeleteByDn(dn, "vnsRsEPgDefToLIf")
}

func (sm *ServiceManager) ReadRelationvnsRsEPgDefToLIf(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsEPgDefToLIf")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsEPgDefToLIf")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvnsRsEPpInfoAtt(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsePpInfoAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsEPpInfoAtt", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvnsRsEPpInfoAtt(parentDn string) error {
	dn := fmt.Sprintf("%s/rsePpInfoAtt", parentDn)
	return sm.DeleteByDn(dn, "vnsRsEPpInfoAtt")
}

func (sm *ServiceManager) ReadRelationvnsRsEPpInfoAtt(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsEPpInfoAtt")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsEPpInfoAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationvnsRsSEPpInfoAtt(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsSEPpInfoAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsSEPpInfoAtt", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvnsRsSEPpInfoAtt(parentDn string) error {
	dn := fmt.Sprintf("%s/rsSEPpInfoAtt", parentDn)
	return sm.DeleteByDn(dn, "vnsRsSEPpInfoAtt")
}

func (sm *ServiceManager) ReadRelationvnsRsSEPpInfoAtt(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsSEPpInfoAtt")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsSEPpInfoAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
