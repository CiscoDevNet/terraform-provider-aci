package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateImportedContract(name string, tenant string, description string, vzCPIfattr models.ImportedContractAttributes) (*models.ImportedContract, error) {
	rn := fmt.Sprintf("cif-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vzCPIf := models.NewImportedContract(rn, parentDn, description, vzCPIfattr)
	err := sm.Save(vzCPIf)
	return vzCPIf, err
}

func (sm *ServiceManager) ReadImportedContract(name string, tenant string) (*models.ImportedContract, error) {
	dn := fmt.Sprintf("uni/tn-%s/cif-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzCPIf := models.ImportedContractFromContainer(cont)
	return vzCPIf, nil
}

func (sm *ServiceManager) DeleteImportedContract(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/cif-%s", tenant, name)
	return sm.DeleteByDn(dn, models.VzcpifClassName)
}

func (sm *ServiceManager) UpdateImportedContract(name string, tenant string, description string, vzCPIfattr models.ImportedContractAttributes) (*models.ImportedContract, error) {
	rn := fmt.Sprintf("cif-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vzCPIf := models.NewImportedContract(rn, parentDn, description, vzCPIfattr)

	vzCPIf.Status = "modified"
	err := sm.Save(vzCPIf)
	return vzCPIf, err

}

func (sm *ServiceManager) ListImportedContract(tenant string) ([]*models.ImportedContract, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/vzCPIf.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ImportedContractListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationvzRsIfFromImportedContract(parentDn, tnVzACtrctName string) error {
	dn := fmt.Sprintf("%s/rsif", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "vzRsIf", dn, tnVzACtrctName))

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

func (sm *ServiceManager) DeleteRelationvzRsIfFromImportedContract(parentDn string) error {
	dn := fmt.Sprintf("%s/rsif", parentDn)
	return sm.DeleteByDn(dn, "vzRsIf")
}

func (sm *ServiceManager) ReadRelationvzRsIfFromImportedContract(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "vzRsIf")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "vzRsIf")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
