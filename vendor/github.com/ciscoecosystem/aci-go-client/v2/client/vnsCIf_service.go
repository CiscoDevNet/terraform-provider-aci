package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateConcreteInterface(name string, parent_dn string, nameAlias string, vnsCIfAttr models.ConcreteInterfaceAttributes) (*models.ConcreteInterface, error) {
	rn := fmt.Sprintf(models.RnvnsCIf, name)
	vnsCIf := models.NewConcreteInterface(rn, parent_dn, nameAlias, vnsCIfAttr)
	err := sm.Save(vnsCIf)
	return vnsCIf, err
}

func (sm *ServiceManager) ReadConcreteInterface(name string, parent_dn string) (*models.ConcreteInterface, error) {
	dn := fmt.Sprintf(models.DnvnsCIf, parent_dn, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsCIf := models.ConcreteInterfaceFromContainer(cont)
	return vnsCIf, nil
}

func (sm *ServiceManager) DeleteConcreteInterface(name string, parent_dn string) error {
	dn := fmt.Sprintf(models.DnvnsCIf, parent_dn, name)
	return sm.DeleteByDn(dn, models.VnscifClassName)
}

func (sm *ServiceManager) UpdateConcreteInterface(name string, parent_dn string, nameAlias string, vnsCIfAttr models.ConcreteInterfaceAttributes) (*models.ConcreteInterface, error) {
	rn := fmt.Sprintf(models.RnvnsCIf, name)
	vnsCIf := models.NewConcreteInterface(rn, parent_dn, nameAlias, vnsCIfAttr)
	vnsCIf.Status = "modified"
	err := sm.Save(vnsCIf)
	return vnsCIf, err
}

func (sm *ServiceManager) ListConcreteInterface(parent_dn string) ([]*models.ConcreteInterface, error) {
	dnUrl := fmt.Sprintf("%s/%s/vnsCIf.json", models.BaseurlStr, parent_dn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ConcreteInterfaceListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsCIfPathAtt(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsCIfPathAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsCIfPathAtt", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationvnsRsCIfPathAtt(parentDn string) error {
	dn := fmt.Sprintf("%s/rsCIfPathAtt", parentDn)
	return sm.DeleteByDn(dn, "vnsRsCIfPathAtt")
}

func (sm *ServiceManager) ReadRelationvnsRsCIfPathAtt(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsCIfPathAtt")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsCIfPathAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
