package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTabooContract(name string, tenant string, description string, vzTabooattr models.TabooContractAttributes) (*models.TabooContract, error) {
	rn := fmt.Sprintf("taboo-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vzTaboo := models.NewTabooContract(rn, parentDn, description, vzTabooattr)
	err := sm.Save(vzTaboo)
	return vzTaboo, err
}

func (sm *ServiceManager) ReadTabooContract(name string, tenant string) (*models.TabooContract, error) {
	dn := fmt.Sprintf("uni/tn-%s/taboo-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzTaboo := models.TabooContractFromContainer(cont)
	return vzTaboo, nil
}

func (sm *ServiceManager) DeleteTabooContract(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/taboo-%s", tenant, name)
	return sm.DeleteByDn(dn, models.VztabooClassName)
}

func (sm *ServiceManager) UpdateTabooContract(name string, tenant string, description string, vzTabooattr models.TabooContractAttributes) (*models.TabooContract, error) {
	rn := fmt.Sprintf("taboo-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vzTaboo := models.NewTabooContract(rn, parentDn, description, vzTabooattr)

	vzTaboo.Status = "modified"
	err := sm.Save(vzTaboo)
	return vzTaboo, err

}

func (sm *ServiceManager) ListTabooContract(tenant string) ([]*models.TabooContract, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/vzTaboo.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.TabooContractListFromContainer(cont)

	return list, err
}
