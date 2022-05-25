package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateConcreteDevice(name string, parent_dn string, nameAlias string, vnsCDevAttr models.ConcreteDeviceAttributes) (*models.ConcreteDevice, error) {
	rn := fmt.Sprintf(models.RnvnsCDev, name)
	vnsCDev := models.NewConcreteDevice(rn, parent_dn, nameAlias, vnsCDevAttr)
	err := sm.Save(vnsCDev)
	return vnsCDev, err
}

func (sm *ServiceManager) ReadConcreteDevice(name string, parent_dn string) (*models.ConcreteDevice, error) {
	dn := fmt.Sprintf(parent_dn+"/"+models.RnvnsCDev, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsCDev := models.ConcreteDeviceFromContainer(cont)
	return vnsCDev, nil
}

func (sm *ServiceManager) DeleteConcreteDevice(name string, parent_dn string) error {
	dn := fmt.Sprintf(parent_dn+"/"+models.RnvnsCDev, name)
	return sm.DeleteByDn(dn, models.VnscdevClassName)
}

func (sm *ServiceManager) UpdateConcreteDevice(name string, parent_dn string, nameAlias string, vnsCDevAttr models.ConcreteDeviceAttributes) (*models.ConcreteDevice, error) {
	rn := fmt.Sprintf(models.RnvnsCDev, name)
	vnsCDev := models.NewConcreteDevice(rn, parent_dn, nameAlias, vnsCDevAttr)
	vnsCDev.Status = "modified"
	err := sm.Save(vnsCDev)
	return vnsCDev, err
}

func (sm *ServiceManager) ListConcreteDevice(parent_dn string) ([]*models.ConcreteDevice, error) {
	dnUrl := fmt.Sprintf(models.BaseurlStr + "/" + parent_dn + "/vnsCDev.json")
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ConcreteDeviceListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsCDevToCtrlrP(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rscDevToCtrlrP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsCDevToCtrlrP", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationvnsRsCDevToCtrlrP(parentDn string) error {
	dn := fmt.Sprintf("%s/rscDevToCtrlrP", parentDn)
	return sm.DeleteByDn(dn, "vnsRsCDevToCtrlrP")
}

func (sm *ServiceManager) ReadRelationvnsRsCDevToCtrlrP(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsCDevToCtrlrP")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsCDevToCtrlrP")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
