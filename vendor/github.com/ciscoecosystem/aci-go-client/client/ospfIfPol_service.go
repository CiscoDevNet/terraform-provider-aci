package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateOSPFInterfacePolicy(name string ,tenant string , description string, ospfIfPolattr models.OSPFInterfacePolicyAttributes) (*models.OSPFInterfacePolicy, error) {	
	rn := fmt.Sprintf("ospfIfPol-%s",name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant )
	ospfIfPol := models.NewOSPFInterfacePolicy(rn, parentDn, description, ospfIfPolattr)
	err := sm.Save(ospfIfPol)
	return ospfIfPol, err
}

func (sm *ServiceManager) ReadOSPFInterfacePolicy(name string ,tenant string ) (*models.OSPFInterfacePolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/ospfIfPol-%s", tenant ,name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfIfPol := models.OSPFInterfacePolicyFromContainer(cont)
	return ospfIfPol, nil
}

func (sm *ServiceManager) DeleteOSPFInterfacePolicy(name string ,tenant string ) error {
	dn := fmt.Sprintf("uni/tn-%s/ospfIfPol-%s", tenant ,name )
	return sm.DeleteByDn(dn, models.OspfifpolClassName)
}

func (sm *ServiceManager) UpdateOSPFInterfacePolicy(name string ,tenant string  ,description string, ospfIfPolattr models.OSPFInterfacePolicyAttributes) (*models.OSPFInterfacePolicy, error) {
	rn := fmt.Sprintf("ospfIfPol-%s",name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant )
	ospfIfPol := models.NewOSPFInterfacePolicy(rn, parentDn, description, ospfIfPolattr)

    ospfIfPol.Status = "modified"
	err := sm.Save(ospfIfPol)
	return ospfIfPol, err

}

func (sm *ServiceManager) ListOSPFInterfacePolicy(tenant string ) ([]*models.OSPFInterfacePolicy, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ospfIfPol.json", baseurlStr , tenant )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.OSPFInterfacePolicyListFromContainer(cont)

	return list, err
}


