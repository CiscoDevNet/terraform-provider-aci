package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outPathAttachment(tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extRsPathL3OutAttattr models.L3outPathAttachmentAttributes) (*models.L3outPathAttachment, error) {
	rn := fmt.Sprintf("rspathL3OutAtt-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	l3extRsPathL3OutAtt := models.NewL3outPathAttachment(rn, parentDn, description, l3extRsPathL3OutAttattr)
	err := sm.Save(l3extRsPathL3OutAtt)
	return l3extRsPathL3OutAtt, err
}

func (sm *ServiceManager) ReadL3outPathAttachment(tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.L3outPathAttachment, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, tDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extRsPathL3OutAtt := models.L3outPathAttachmentFromContainer(cont)
	return l3extRsPathL3OutAtt, nil
}

func (sm *ServiceManager) DeleteL3outPathAttachment(tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/rspathL3OutAtt-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, tDn)
	return sm.DeleteByDn(dn, models.L3extrspathl3outattClassName)
}

func (sm *ServiceManager) UpdateL3outPathAttachment(tDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extRsPathL3OutAttattr models.L3outPathAttachmentAttributes) (*models.L3outPathAttachment, error) {
	rn := fmt.Sprintf("rspathL3OutAtt-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	l3extRsPathL3OutAtt := models.NewL3outPathAttachment(rn, parentDn, description, l3extRsPathL3OutAttattr)

	l3extRsPathL3OutAtt.Status = "modified"
	err := sm.Save(l3extRsPathL3OutAtt)
	return l3extRsPathL3OutAtt, err

}

func (sm *ServiceManager) ListL3outPathAttachment(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.L3outPathAttachment, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/l3extRsPathL3OutAtt.json", baseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outPathAttachmentListFromContainer(cont)

	return list, err
}
