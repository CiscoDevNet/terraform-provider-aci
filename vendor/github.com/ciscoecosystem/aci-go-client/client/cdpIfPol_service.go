package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCDPInterfacePolicy(name string, description string, cdpIfPolattr models.CDPInterfacePolicyAttributes) (*models.CDPInterfacePolicy, error) {
	rn := fmt.Sprintf("infra/cdpIfP-%s", name)
	parentDn := fmt.Sprintf("uni")
	cdpIfPol := models.NewCDPInterfacePolicy(rn, parentDn, description, cdpIfPolattr)
	err := sm.Save(cdpIfPol)
	return cdpIfPol, err
}

func (sm *ServiceManager) ReadCDPInterfacePolicy(name string) (*models.CDPInterfacePolicy, error) {
	dn := fmt.Sprintf("uni/infra/cdpIfP-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cdpIfPol := models.CDPInterfacePolicyFromContainer(cont)
	return cdpIfPol, nil
}

func (sm *ServiceManager) DeleteCDPInterfacePolicy(name string) error {
	dn := fmt.Sprintf("uni/infra/cdpIfP-%s", name)
	return sm.DeleteByDn(dn, models.CdpifpolClassName)
}

func (sm *ServiceManager) UpdateCDPInterfacePolicy(name string, description string, cdpIfPolattr models.CDPInterfacePolicyAttributes) (*models.CDPInterfacePolicy, error) {
	rn := fmt.Sprintf("infra/cdpIfP-%s", name)
	parentDn := fmt.Sprintf("uni")
	cdpIfPol := models.NewCDPInterfacePolicy(rn, parentDn, description, cdpIfPolattr)

	cdpIfPol.Status = "modified"
	err := sm.Save(cdpIfPol)
	return cdpIfPol, err

}

func (sm *ServiceManager) ListCDPInterfacePolicy() ([]*models.CDPInterfacePolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/cdpIfPol.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.CDPInterfacePolicyListFromContainer(cont)

	return list, err
}
