package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outVPCMember(side string, leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extMemberattr models.L3outVPCMemberAttributes) (*models.L3outVPCMember, error) {
	rn := fmt.Sprintf("mem-%s", side)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn)
	l3extMember := models.NewL3outVPCMember(rn, parentDn, description, l3extMemberattr)
	err := sm.Save(l3extMember)
	return l3extMember, err
}

func (sm *ServiceManager) ReadL3outVPCMember(side string, leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.L3outVPCMember, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]/mem-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn, side)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extMember := models.L3outVPCMemberFromContainer(cont)
	return l3extMember, nil
}

func (sm *ServiceManager) DeleteL3outVPCMember(side string, leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]/mem-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn, side)
	return sm.DeleteByDn(dn, models.L3extmemberClassName)
}

func (sm *ServiceManager) UpdateL3outVPCMember(side string, leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extMemberattr models.L3outVPCMemberAttributes) (*models.L3outVPCMember, error) {
	rn := fmt.Sprintf("mem-%s", side)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn)
	l3extMember := models.NewL3outVPCMember(rn, parentDn, description, l3extMemberattr)

	l3extMember.Status = "modified"
	err := sm.Save(l3extMember)
	return l3extMember, err

}

func (sm *ServiceManager) ListL3outVPCMember(leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.L3outVPCMember, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]/l3extMember.json", baseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outVPCMemberListFromContainer(cont)

	return list, err
}
