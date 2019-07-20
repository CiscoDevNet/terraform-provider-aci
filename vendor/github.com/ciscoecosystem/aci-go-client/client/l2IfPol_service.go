package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateL2InterfacePolicy(name string , description string, l2IfPolattr models.L2InterfacePolicyAttributes) (*models.L2InterfacePolicy, error) {	
	rn := fmt.Sprintf("infra/l2IfP-%s",name)
	parentDn := fmt.Sprintf("uni")
	l2IfPol := models.NewL2InterfacePolicy(rn, parentDn, description, l2IfPolattr)
	err := sm.Save(l2IfPol)
	return l2IfPol, err
}

func (sm *ServiceManager) ReadL2InterfacePolicy(name string ) (*models.L2InterfacePolicy, error) {
	dn := fmt.Sprintf("uni/infra/l2IfP-%s", name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l2IfPol := models.L2InterfacePolicyFromContainer(cont)
	return l2IfPol, nil
}

func (sm *ServiceManager) DeleteL2InterfacePolicy(name string ) error {
	dn := fmt.Sprintf("uni/infra/l2IfP-%s", name )
	return sm.DeleteByDn(dn, models.L2ifpolClassName)
}

func (sm *ServiceManager) UpdateL2InterfacePolicy(name string  ,description string, l2IfPolattr models.L2InterfacePolicyAttributes) (*models.L2InterfacePolicy, error) {
	rn := fmt.Sprintf("infra/l2IfP-%s",name)
	parentDn := fmt.Sprintf("uni")
	l2IfPol := models.NewL2InterfacePolicy(rn, parentDn, description, l2IfPolattr)

    l2IfPol.Status = "modified"
	err := sm.Save(l2IfPol)
	return l2IfPol, err

}

func (sm *ServiceManager) ListL2InterfacePolicy() ([]*models.L2InterfacePolicy, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/l2IfPol.json", baseurlStr )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.L2InterfacePolicyListFromContainer(cont)

	return list, err
}


