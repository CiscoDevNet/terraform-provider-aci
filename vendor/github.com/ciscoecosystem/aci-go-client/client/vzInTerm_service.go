package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateInTermSubject(contract_subject string, contract string, tenant string, description string, nameAlias string, vzInTermAttr models.InTermSubjectAttributes) (*models.InTermSubject, error) {
	rn := fmt.Sprintf(models.RnvzInTerm)
	parentDn := fmt.Sprintf(models.ParentDnvzInTerm, tenant, contract, contract_subject)
	vzInTerm := models.NewInTermSubject(rn, parentDn, description, nameAlias, vzInTermAttr)
	err := sm.Save(vzInTerm)
	return vzInTerm, err
}

func (sm *ServiceManager) ReadInTermSubject(contract_subject string, contract string, tenant string) (*models.InTermSubject, error) {
	dn := fmt.Sprintf(models.DnvzInTerm, tenant, contract, contract_subject)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzInTerm := models.InTermSubjectFromContainer(cont)
	return vzInTerm, nil
}

func (sm *ServiceManager) DeleteInTermSubject(contract_subject string, contract string, tenant string) error {
	dn := fmt.Sprintf(models.DnvzInTerm, tenant, contract, contract_subject)
	return sm.DeleteByDn(dn, models.VzintermClassName)
}

func (sm *ServiceManager) UpdateInTermSubject(contract_subject string, contract string, tenant string, description string, nameAlias string, vzInTermAttr models.InTermSubjectAttributes) (*models.InTermSubject, error) {
	rn := fmt.Sprintf(models.RnvzInTerm)
	parentDn := fmt.Sprintf(models.ParentDnvzInTerm, tenant, contract, contract_subject)
	vzInTerm := models.NewInTermSubject(rn, parentDn, description, nameAlias, vzInTermAttr)
	vzInTerm.Status = "modified"
	err := sm.Save(vzInTerm)
	return vzInTerm, err
}

func (sm *ServiceManager) ListInTermSubject(contract_subject string, contract string, tenant string) ([]*models.InTermSubject, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/brc-%s/subj-%s/vzInTerm.json", models.BaseurlStr, tenant, contract, contract_subject)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.InTermSubjectListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvzRsInTermGraphAtt(parentDn, annotation, tnVnsAbsGraphName string) error {
	dn := fmt.Sprintf("%s/rsInTermGraphAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnVnsAbsGraphName": "%s"
			}
		}
	}`, "vzRsInTermGraphAtt", dn, annotation, tnVnsAbsGraphName))

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

func (sm *ServiceManager) DeleteRelationvzRsInTermGraphAtt(parentDn string) error {
	dn := fmt.Sprintf("%s/rsInTermGraphAtt", parentDn)
	return sm.DeleteByDn(dn, "vzRsInTermGraphAtt")
}

func (sm *ServiceManager) ReadRelationvzRsInTermGraphAtt(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vzRsInTermGraphAtt")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vzRsInTermGraphAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}
