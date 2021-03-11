package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outPathAttachmentSecondaryIp(addr string, leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extIpattr models.L3outPathAttachmentSecondaryIpAttributes) (*models.L3outPathAttachmentSecondaryIp, error) {
	rn := fmt.Sprintf("addr-[%s]", addr)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn)
	l3extIp := models.NewL3outPathAttachmentSecondaryIp(rn, parentDn, description, l3extIpattr)
	err := sm.Save(l3extIp)
	return l3extIp, err
}

func (sm *ServiceManager) ReadL3outPathAttachmentSecondaryIp(addr string, leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.L3outPathAttachmentSecondaryIp, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]/addr-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn, addr)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extIp := models.L3outPathAttachmentSecondaryIpFromContainer(cont)
	return l3extIp, nil
}

func (sm *ServiceManager) DeleteL3outPathAttachmentSecondaryIp(addr string, leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]/addr-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn, addr)
	return sm.DeleteByDn(dn, models.L3extipClassName)
}

func (sm *ServiceManager) UpdateL3outPathAttachmentSecondaryIp(addr string, leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extIpattr models.L3outPathAttachmentSecondaryIpAttributes) (*models.L3outPathAttachmentSecondaryIp, error) {
	rn := fmt.Sprintf("addr-[%s]", addr)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn)
	l3extIp := models.NewL3outPathAttachmentSecondaryIp(rn, parentDn, description, l3extIpattr)

	l3extIp.Status = "modified"
	err := sm.Save(l3extIp)
	return l3extIp, err

}

func (sm *ServiceManager) ListL3outPathAttachmentSecondaryIp(leaf_port_tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.L3outPathAttachmentSecondaryIp, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]/l3extIp.json", baseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile, leaf_port_tDn)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outPathAttachmentSecondaryIpListFromContainer(cont)

	return list, err
}
