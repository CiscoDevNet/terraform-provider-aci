package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTACACSSource(name string, monitoring_target_scope string, monitoring_policy string, tenant string, description string, nameAlias string, tacacsSrcAttr models.TACACSSourceAttributes) (*models.TACACSSource, error) {
	rn := fmt.Sprintf(models.RntacacsSrc, name)
	parentDn := fmt.Sprintf(models.ParentDntacacsSrc, tenant, monitoring_policy, monitoring_target_scope, monitoring_policy)
	tacacsSrc := models.NewTACACSSource(rn, parentDn, description, nameAlias, tacacsSrcAttr)
	err := sm.Save(tacacsSrc)
	return tacacsSrc, err
}

func (sm *ServiceManager) ReadTACACSSource(name string, monitoring_target_scope string, monitoring_policy string, tenant string) (*models.TACACSSource, error) {
	dn := fmt.Sprintf(models.DntacacsSrc, tenant, monitoring_policy, monitoring_target_scope, monitoring_policy, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	tacacsSrc := models.TACACSSourceFromContainer(cont)
	return tacacsSrc, nil
}

func (sm *ServiceManager) DeleteTACACSSource(name string, monitoring_target_scope string, monitoring_policy string, tenant string) error {
	dn := fmt.Sprintf(models.DntacacsSrc, tenant, monitoring_policy, monitoring_target_scope, monitoring_policy, name)
	return sm.DeleteByDn(dn, models.TacacssrcClassName)
}

func (sm *ServiceManager) UpdateTACACSSource(name string, monitoring_target_scope string, monitoring_policy string, tenant string, description string, nameAlias string, tacacsSrcAttr models.TACACSSourceAttributes) (*models.TACACSSource, error) {
	rn := fmt.Sprintf(models.RntacacsSrc, name)
	parentDn := fmt.Sprintf(models.ParentDntacacsSrc, tenant, monitoring_policy, monitoring_target_scope, monitoring_policy)
	tacacsSrc := models.NewTACACSSource(rn, parentDn, description, nameAlias, tacacsSrcAttr)
	tacacsSrc.Status = "modified"
	err := sm.Save(tacacsSrc)
	return tacacsSrc, err
}

func (sm *ServiceManager) ListTACACSSource(monitoring_target_scope string, monitoring_policy string, tenant string) ([]*models.TACACSSource, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/monepg-%s/tarepg-%s/tacacsSrc.json", models.BaseurlStr, tenant, monitoring_policy, monitoring_target_scope, monitoring_policy)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.TACACSSourceListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationtacacsRsDestGroup(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsdestGroup", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "tacacsRsDestGroup", dn, annotation, tDn))

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

func (sm *ServiceManager) DeleteRelationtacacsRsDestGroup(parentDn string) error {
	dn := fmt.Sprintf("%s/rsdestGroup", parentDn)
	return sm.DeleteByDn(dn, "tacacsRsDestGroup")
}

func (sm *ServiceManager) ReadRelationtacacsRsDestGroup(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "tacacsRsDestGroup")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "tacacsRsDestGroup")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
