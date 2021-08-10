package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateContractSubject(name string, contract string, tenant string, description string, vzSubjattr models.ContractSubjectAttributes) (*models.ContractSubject, error) {
	rn := fmt.Sprintf("subj-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/brc-%s", tenant, contract)
	vzSubj := models.NewContractSubject(rn, parentDn, description, vzSubjattr)
	err := sm.Save(vzSubj)
	return vzSubj, err
}

func (sm *ServiceManager) ReadContractSubject(name string, contract string, tenant string) (*models.ContractSubject, error) {
	dn := fmt.Sprintf("uni/tn-%s/brc-%s/subj-%s", tenant, contract, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzSubj := models.ContractSubjectFromContainer(cont)
	return vzSubj, nil
}

func (sm *ServiceManager) DeleteContractSubject(name string, contract string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/brc-%s/subj-%s", tenant, contract, name)
	return sm.DeleteByDn(dn, models.VzsubjClassName)
}

func (sm *ServiceManager) UpdateContractSubject(name string, contract string, tenant string, description string, vzSubjattr models.ContractSubjectAttributes) (*models.ContractSubject, error) {
	rn := fmt.Sprintf("subj-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/brc-%s", tenant, contract)
	vzSubj := models.NewContractSubject(rn, parentDn, description, vzSubjattr)

	vzSubj.Status = "modified"
	err := sm.Save(vzSubj)
	return vzSubj, err

}

func (sm *ServiceManager) ListContractSubject(contract string, tenant string) ([]*models.ContractSubject, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/brc-%s/vzSubj.json", baseurlStr, tenant, contract)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ContractSubjectListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvzRsSubjGraphAttFromContractSubject(parentDn, tnVnsAbsGraphName string) error {
	dn := fmt.Sprintf("%s/rsSubjGraphAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVnsAbsGraphName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vzRsSubjGraphAtt", dn, tnVnsAbsGraphName))

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

func (sm *ServiceManager) DeleteRelationvzRsSubjGraphAttFromContractSubject(parentDn string) error {
	dn := fmt.Sprintf("%s/rsSubjGraphAtt", parentDn)
	return sm.DeleteByDn(dn, "vzRsSubjGraphAtt")
}

func (sm *ServiceManager) ReadRelationvzRsSubjGraphAttFromContractSubject(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsSubjGraphAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsSubjGraphAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvzRsSdwanPolFromContractSubject(parentDn, tnExtdevSDWanSlaPolName string) error {
	dn := fmt.Sprintf("%s/rsSdwanPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vzRsSdwanPol", dn, tnExtdevSDWanSlaPolName))

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

func (sm *ServiceManager) DeleteRelationvzRsSdwanPolFromContractSubject(parentDn string) error {
	dn := fmt.Sprintf("%s/rsSdwanPol", parentDn)
	return sm.DeleteByDn(dn, "vzRsSdwanPol")
}

func (sm *ServiceManager) ReadRelationvzRsSdwanPolFromContractSubject(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsSdwanPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsSdwanPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationvzRsSubjFiltAttFromContractSubject(parentDn, tnVzFilterName string) error {
	dn := fmt.Sprintf("%s/rssubjFiltAtt-%s", parentDn, tnVzFilterName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "vzRsSubjFiltAtt", dn))

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

func (sm *ServiceManager) DeleteRelationvzRsSubjFiltAttFromContractSubject(parentDn, tnVzFilterName string) error {
	dn := fmt.Sprintf("%s/rssubjFiltAtt-%s", parentDn, tnVzFilterName)
	return sm.DeleteByDn(dn, "vzRsSubjFiltAtt")
}

func (sm *ServiceManager) ReadRelationvzRsSubjFiltAttFromContractSubject(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsSubjFiltAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsSubjFiltAtt")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
