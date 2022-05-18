package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateOutTermSubject(contract_subject string, contract string, tenant string, description string, nameAlias string, vzOutTermAttr models.OutTermSubjectAttributes) (*models.OutTermSubject, error) {
	rn := fmt.Sprintf(models.RnvzOutTerm)
	parentDn := fmt.Sprintf(models.ParentDnvzOutTerm, tenant, contract, contract_subject)
	vzOutTerm := models.NewOutTermSubject(rn, parentDn, description, nameAlias, vzOutTermAttr)
	err := sm.Save(vzOutTerm)
	return vzOutTerm, err
}

func (sm *ServiceManager) ReadOutTermSubject(contract_subject string, contract string, tenant string) (*models.OutTermSubject, error) {
	dn := fmt.Sprintf(models.DnvzOutTerm, tenant, contract, contract_subject)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzOutTerm := models.OutTermSubjectFromContainer(cont)
	return vzOutTerm, nil
}

func (sm *ServiceManager) DeleteOutTermSubject(contract_subject string, contract string, tenant string) error {
	dn := fmt.Sprintf(models.DnvzOutTerm, tenant, contract, contract_subject)
	return sm.DeleteByDn(dn, models.VzouttermClassName)
}

func (sm *ServiceManager) UpdateOutTermSubject(contract_subject string, contract string, tenant string, description string, nameAlias string, vzOutTermAttr models.OutTermSubjectAttributes) (*models.OutTermSubject, error) {
	rn := fmt.Sprintf(models.RnvzOutTerm)
	parentDn := fmt.Sprintf(models.ParentDnvzOutTerm, tenant, contract, contract_subject)
	vzOutTerm := models.NewOutTermSubject(rn, parentDn, description, nameAlias, vzOutTermAttr)
	vzOutTerm.Status = "modified"
	err := sm.Save(vzOutTerm)
	return vzOutTerm, err
}

func (sm *ServiceManager) ListOutTermSubject(contract_subject string, contract string, tenant string) ([]*models.OutTermSubject, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/brc-%s/subj-%s/vzOutTerm.json", models.BaseurlStr, tenant, contract, contract_subject)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.OutTermSubjectListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvzRsOutTermGraphAtt(parentDn, annotation, tnVnsAbsGraphName string) error {
	dn := fmt.Sprintf("%s/rsOutTermGraphAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnVnsAbsGraphName": "%s"
			}
		}
	}`, "vzRsOutTermGraphAtt", dn, annotation, tnVnsAbsGraphName))

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

func (sm *ServiceManager) DeleteRelationvzRsOutTermGraphAtt(parentDn string) error {
	dn := fmt.Sprintf("%s/rsOutTermGraphAtt", parentDn)
	return sm.DeleteByDn(dn, "vzRsOutTermGraphAtt")
}

func (sm *ServiceManager) ReadRelationvzRsOutTermGraphAtt(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vzRsOutTermGraphAtt")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vzRsOutTermGraphAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
