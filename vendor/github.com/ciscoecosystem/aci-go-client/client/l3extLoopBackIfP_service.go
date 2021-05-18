package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLoopBackInterfaceProfile(addr string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string, description string, l3extLoopBackIfPattr models.LoopBackInterfaceProfileAttributes) (*models.LoopBackInterfaceProfile, error) {
	rn := fmt.Sprintf("lbp-[%s]", addr)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn)
	l3extLoopBackIfP := models.NewLoopBackInterfaceProfile(rn, parentDn, description, l3extLoopBackIfPattr)
	err := sm.Save(l3extLoopBackIfP)
	return l3extLoopBackIfP, err
}

func (sm *ServiceManager) ReadLoopBackInterfaceProfile(addr string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string) (*models.LoopBackInterfaceProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/lbp-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn, addr)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extLoopBackIfP := models.LoopBackInterfaceProfileFromContainer(cont)
	return l3extLoopBackIfP, nil
}

func (sm *ServiceManager) DeleteLoopBackInterfaceProfile(addr string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/lbp-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn, addr)
	return sm.DeleteByDn(dn, models.L3extloopbackifpClassName)
}

func (sm *ServiceManager) UpdateLoopBackInterfaceProfile(addr string, fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string, description string, l3extLoopBackIfPattr models.LoopBackInterfaceProfileAttributes) (*models.LoopBackInterfaceProfile, error) {
	rn := fmt.Sprintf("lbp-[%s]", addr)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, fabric_node_tDn)
	l3extLoopBackIfP := models.NewLoopBackInterfaceProfile(rn, parentDn, description, l3extLoopBackIfPattr)

	l3extLoopBackIfP.Status = "modified"
	err := sm.Save(l3extLoopBackIfP)
	return l3extLoopBackIfP, err

}

func (sm *ServiceManager) ListLoopBackInterfaceProfile(fabric_node_tDn string, logical_node_profile string, l3_outside string, tenant string) ([]*models.LoopBackInterfaceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/l3extLoopBackIfP.json", baseurlStr, tenant, l3_outside, logical_node_profile, fabric_node_tDn)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LoopBackInterfaceProfileListFromContainer(cont)

	return list, err
}
