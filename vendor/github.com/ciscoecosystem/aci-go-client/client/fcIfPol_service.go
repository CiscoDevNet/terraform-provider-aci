package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateInterfaceFCPolicy(name string , description string, fcIfPolattr models.InterfaceFCPolicyAttributes) (*models.InterfaceFCPolicy, error) {	
	rn := fmt.Sprintf("infra/fcIfPol-%s",name)
	parentDn := fmt.Sprintf("uni")
	fcIfPol := models.NewInterfaceFCPolicy(rn, parentDn, description, fcIfPolattr)
	err := sm.Save(fcIfPol)
	return fcIfPol, err
}

func (sm *ServiceManager) ReadInterfaceFCPolicy(name string ) (*models.InterfaceFCPolicy, error) {
	dn := fmt.Sprintf("uni/infra/fcIfPol-%s", name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fcIfPol := models.InterfaceFCPolicyFromContainer(cont)
	return fcIfPol, nil
}

func (sm *ServiceManager) DeleteInterfaceFCPolicy(name string ) error {
	dn := fmt.Sprintf("uni/infra/fcIfPol-%s", name )
	return sm.DeleteByDn(dn, models.FcifpolClassName)
}

func (sm *ServiceManager) UpdateInterfaceFCPolicy(name string  ,description string, fcIfPolattr models.InterfaceFCPolicyAttributes) (*models.InterfaceFCPolicy, error) {
	rn := fmt.Sprintf("infra/fcIfPol-%s",name)
	parentDn := fmt.Sprintf("uni")
	fcIfPol := models.NewInterfaceFCPolicy(rn, parentDn, description, fcIfPolattr)

    fcIfPol.Status = "modified"
	err := sm.Save(fcIfPol)
	return fcIfPol, err

}

func (sm *ServiceManager) ListInterfaceFCPolicy() ([]*models.InterfaceFCPolicy, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/fcIfPol.json", baseurlStr )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.InterfaceFCPolicyListFromContainer(cont)

	return list, err
}


