package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLLDPInterfacePolicy(name string, description string, lldpIfPolattr models.LLDPInterfacePolicyAttributes) (*models.LLDPInterfacePolicy, error) {
	rn := fmt.Sprintf("infra/lldpIfP-%s", name)
	parentDn := fmt.Sprintf("uni")
	lldpIfPol := models.NewLLDPInterfacePolicy(rn, parentDn, description, lldpIfPolattr)
	err := sm.Save(lldpIfPol)
	return lldpIfPol, err
}

func (sm *ServiceManager) ReadLLDPInterfacePolicy(name string) (*models.LLDPInterfacePolicy, error) {
	dn := fmt.Sprintf("uni/infra/lldpIfP-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	lldpIfPol := models.LLDPInterfacePolicyFromContainer(cont)
	return lldpIfPol, nil
}

func (sm *ServiceManager) DeleteLLDPInterfacePolicy(name string) error {
	dn := fmt.Sprintf("uni/infra/lldpIfP-%s", name)
	return sm.DeleteByDn(dn, models.LldpifpolClassName)
}

func (sm *ServiceManager) UpdateLLDPInterfacePolicy(name string, description string, lldpIfPolattr models.LLDPInterfacePolicyAttributes) (*models.LLDPInterfacePolicy, error) {
	rn := fmt.Sprintf("infra/lldpIfP-%s", name)
	parentDn := fmt.Sprintf("uni")
	lldpIfPol := models.NewLLDPInterfacePolicy(rn, parentDn, description, lldpIfPolattr)

	lldpIfPol.Status = "modified"
	err := sm.Save(lldpIfPol)
	return lldpIfPol, err

}

func (sm *ServiceManager) ListLLDPInterfacePolicy() ([]*models.LLDPInterfacePolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/lldpIfPol.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LLDPInterfacePolicyListFromContainer(cont)

	return list, err
}
