package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateLACPPolicy(name string , description string, lacpLagPolattr models.LACPPolicyAttributes) (*models.LACPPolicy, error) {	
	rn := fmt.Sprintf("infra/lacplagp-%s",name)
	parentDn := fmt.Sprintf("uni")
	lacpLagPol := models.NewLACPPolicy(rn, parentDn, description, lacpLagPolattr)
	err := sm.Save(lacpLagPol)
	return lacpLagPol, err
}

func (sm *ServiceManager) ReadLACPPolicy(name string ) (*models.LACPPolicy, error) {
	dn := fmt.Sprintf("uni/infra/lacplagp-%s", name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	lacpLagPol := models.LACPPolicyFromContainer(cont)
	return lacpLagPol, nil
}

func (sm *ServiceManager) DeleteLACPPolicy(name string ) error {
	dn := fmt.Sprintf("uni/infra/lacplagp-%s", name )
	return sm.DeleteByDn(dn, models.LacplagpolClassName)
}

func (sm *ServiceManager) UpdateLACPPolicy(name string  ,description string, lacpLagPolattr models.LACPPolicyAttributes) (*models.LACPPolicy, error) {
	rn := fmt.Sprintf("infra/lacplagp-%s",name)
	parentDn := fmt.Sprintf("uni")
	lacpLagPol := models.NewLACPPolicy(rn, parentDn, description, lacpLagPolattr)

    lacpLagPol.Status = "modified"
	err := sm.Save(lacpLagPol)
	return lacpLagPol, err

}

func (sm *ServiceManager) ListLACPPolicy() ([]*models.LACPPolicy, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/lacpLagPol.json", baseurlStr )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.LACPPolicyListFromContainer(cont)

	return list, err
}


