package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateRelationfromaAbsNodetoanLDev(function_node string ,l4-l7_service_graph_template string ,tenant string , description string, vnsRsNodeToLDevattr models.RelationfromaAbsNodetoanLDevAttributes) (*models.RelationfromaAbsNodetoanLDev, error) {	
	rn := fmt.Sprintf("rsNodeToLDev")
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s", tenant ,l4-l7_service_graph_template ,function_node )
	vnsRsNodeToLDev := models.NewRelationfromaAbsNodetoanLDev(rn, parentDn, description, vnsRsNodeToLDevattr)
	err := sm.Save(vnsRsNodeToLDev)
	return vnsRsNodeToLDev, err
}

func (sm *ServiceManager) ReadRelationfromaAbsNodetoanLDev(function_node string ,l4-l7_service_graph_template string ,tenant string ) (*models.RelationfromaAbsNodetoanLDev, error) {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s/rsNodeToLDev", tenant ,l4-l7_service_graph_template ,function_node )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsRsNodeToLDev := models.RelationfromaAbsNodetoanLDevFromContainer(cont)
	return vnsRsNodeToLDev, nil
}

func (sm *ServiceManager) DeleteRelationfromaAbsNodetoanLDev(function_node string ,l4-l7_service_graph_template string ,tenant string ) error {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s/rsNodeToLDev", tenant ,l4-l7_service_graph_template ,function_node )
	return sm.DeleteByDn(dn, models.VnsrsnodetoldevClassName)
}

func (sm *ServiceManager) UpdateRelationfromaAbsNodetoanLDev(function_node string ,l4-l7_service_graph_template string ,tenant string  ,description string, vnsRsNodeToLDevattr models.RelationfromaAbsNodetoanLDevAttributes) (*models.RelationfromaAbsNodetoanLDev, error) {
	rn := fmt.Sprintf("rsNodeToLDev")
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsNode-%s", tenant ,l4-l7_service_graph_template ,function_node )
	vnsRsNodeToLDev := models.NewRelationfromaAbsNodetoanLDev(rn, parentDn, description, vnsRsNodeToLDevattr)

    vnsRsNodeToLDev.Status = "modified"
	err := sm.Save(vnsRsNodeToLDev)
	return vnsRsNodeToLDev, err

}

func (sm *ServiceManager) ListRelationfromaAbsNodetoanLDev(function_node string ,l4-l7_service_graph_template string ,tenant string ) ([]*models.RelationfromaAbsNodetoanLDev, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/AbsGraph-%s/AbsNode-%s/vnsRsNodeToLDev.json", baseurlStr , tenant ,l4-l7_service_graph_template ,function_node )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.RelationfromaAbsNodetoanLDevListFromContainer(cont)

	return list, err
}


