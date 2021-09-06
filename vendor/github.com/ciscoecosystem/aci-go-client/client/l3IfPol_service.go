package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3InterfacePolicy(name string, description string, l3IfPolattr models.L3InterfacePolicyAttributes) (*models.L3InterfacePolicy, error) {
	rn := fmt.Sprintf("fabric/l3IfP-%s", name)
	parentDn := fmt.Sprintf("uni")
	l3IfPol := models.NewL3InterfacePolicy(rn, parentDn, description, l3IfPolattr)
	err := sm.Save(l3IfPol)
	return l3IfPol, err
}

func (sm *ServiceManager) ReadL3InterfacePolicy(name string) (*models.L3InterfacePolicy, error) {
	dn := fmt.Sprintf("uni/fabric/l3IfP-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3IfPol := models.L3InterfacePolicyFromContainer(cont)
	return l3IfPol, nil
}

func (sm *ServiceManager) DeleteL3InterfacePolicy(name string) error {
	dn := fmt.Sprintf("uni/fabric/l3IfP-%s", name)
	return sm.DeleteByDn(dn, models.L3ifpolClassName)
}

func (sm *ServiceManager) UpdateL3InterfacePolicy(name string, description string, l3IfPolattr models.L3InterfacePolicyAttributes) (*models.L3InterfacePolicy, error) {
	rn := fmt.Sprintf("fabric/l3IfP-%s", name)
	parentDn := fmt.Sprintf("uni")
	l3IfPol := models.NewL3InterfacePolicy(rn, parentDn, description, l3IfPolattr)

	l3IfPol.Status = "modified"
	err := sm.Save(l3IfPol)
	return l3IfPol, err

}

func (sm *ServiceManager) ListL3InterfacePolicy() ([]*models.L3InterfacePolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/l3IfPol.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3InterfacePolicyListFromContainer(cont)

	return list, err
}
