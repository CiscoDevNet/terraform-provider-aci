package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBFDInterfacePolicy(name string, tenant string, description string, bfdIfPolattr models.BFDInterfacePolicyAttributes) (*models.BFDInterfacePolicy, error) {
	rn := fmt.Sprintf("bfdIfPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bfdIfPol := models.NewBFDInterfacePolicy(rn, parentDn, description, bfdIfPolattr)
	err := sm.Save(bfdIfPol)
	return bfdIfPol, err
}

func (sm *ServiceManager) ReadBFDInterfacePolicy(name string, tenant string) (*models.BFDInterfacePolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/bfdIfPol-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	bfdIfPol := models.BFDInterfacePolicyFromContainer(cont)
	return bfdIfPol, nil
}

func (sm *ServiceManager) DeleteBFDInterfacePolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/bfdIfPol-%s", tenant, name)
	return sm.DeleteByDn(dn, models.BfdifpolClassName)
}

func (sm *ServiceManager) UpdateBFDInterfacePolicy(name string, tenant string, description string, bfdIfPolattr models.BFDInterfacePolicyAttributes) (*models.BFDInterfacePolicy, error) {
	rn := fmt.Sprintf("bfdIfPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	bfdIfPol := models.NewBFDInterfacePolicy(rn, parentDn, description, bfdIfPolattr)

	bfdIfPol.Status = "modified"
	err := sm.Save(bfdIfPol)
	return bfdIfPol, err

}

func (sm *ServiceManager) ListBFDInterfacePolicy(tenant string) ([]*models.BFDInterfacePolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/bfdIfPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.BFDInterfacePolicyListFromContainer(cont)

	return list, err
}
