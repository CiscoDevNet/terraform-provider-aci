package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateContract(name string, tenant string, description string, vzBrCPattr models.ContractAttributes) (*models.Contract, error) {
	rn := fmt.Sprintf("brc-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vzBrCP := models.NewContract(rn, parentDn, description, vzBrCPattr)
	err := sm.Save(vzBrCP)
	return vzBrCP, err
}

func (sm *ServiceManager) ReadContract(name string, tenant string) (*models.Contract, error) {
	dn := fmt.Sprintf("uni/tn-%s/brc-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzBrCP := models.ContractFromContainer(cont)
	return vzBrCP, nil
}

func (sm *ServiceManager) DeleteContract(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/brc-%s", tenant, name)
	return sm.DeleteByDn(dn, models.VzbrcpClassName)
}

func (sm *ServiceManager) UpdateContract(name string, tenant string, description string, vzBrCPattr models.ContractAttributes) (*models.Contract, error) {
	rn := fmt.Sprintf("brc-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vzBrCP := models.NewContract(rn, parentDn, description, vzBrCPattr)

	vzBrCP.Status = "modified"
	err := sm.Save(vzBrCP)
	return vzBrCP, err

}

func (sm *ServiceManager) ListContract(tenant string) ([]*models.Contract, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/vzBrCP.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ContractListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvzRsGraphAttFromContract(parentDn, tnVnsAbsGraphName string) error {
	dn := fmt.Sprintf("%s/rsGraphAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVnsAbsGraphName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "vzRsGraphAtt", dn, tnVnsAbsGraphName))

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

func (sm *ServiceManager) DeleteRelationvzRsGraphAttFromContract(parentDn string) error {
	dn := fmt.Sprintf("%s/rsGraphAtt", parentDn)
	return sm.DeleteByDn(dn, "vzRsGraphAtt")
}

func (sm *ServiceManager) ReadRelationvzRsGraphAttFromContract(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsGraphAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsGraphAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
