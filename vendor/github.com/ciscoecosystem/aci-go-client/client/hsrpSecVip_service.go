package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL3outHSRPSecondaryVIP(ip string, hsrp_group_profile string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, hsrpSecVipattr models.L3outHSRPSecondaryVIPAttributes) (*models.L3outHSRPSecondaryVIP, error) {
	rn := fmt.Sprintf("hsrpSecVip-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile, hsrp_group_profile)
	hsrpSecVip := models.NewL3outHSRPSecondaryVIP(rn, parentDn, description, hsrpSecVipattr)
	err := sm.Save(hsrpSecVip)
	return hsrpSecVip, err
}

func (sm *ServiceManager) ReadL3outHSRPSecondaryVIP(ip string, hsrp_group_profile string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.L3outHSRPSecondaryVIP, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s/hsrpSecVip-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, hsrp_group_profile, ip)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpSecVip := models.L3outHSRPSecondaryVIPFromContainer(cont)
	return hsrpSecVip, nil
}

func (sm *ServiceManager) DeleteL3outHSRPSecondaryVIP(ip string, hsrp_group_profile string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s/hsrpSecVip-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, hsrp_group_profile, ip)
	return sm.DeleteByDn(dn, models.HsrpsecvipClassName)
}

func (sm *ServiceManager) UpdateL3outHSRPSecondaryVIP(ip string, hsrp_group_profile string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, hsrpSecVipattr models.L3outHSRPSecondaryVIPAttributes) (*models.L3outHSRPSecondaryVIP, error) {
	rn := fmt.Sprintf("hsrpSecVip-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile, hsrp_group_profile)
	hsrpSecVip := models.NewL3outHSRPSecondaryVIP(rn, parentDn, description, hsrpSecVipattr)

	hsrpSecVip.Status = "modified"
	err := sm.Save(hsrpSecVip)
	return hsrpSecVip, err

}

func (sm *ServiceManager) ListL3outHSRPSecondaryVIP(hsrp_group_profile string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.L3outHSRPSecondaryVIP, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/hsrpIfP/hsrpGroupP-%s/hsrpSecVip.json", baseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile, hsrp_group_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3outHSRPSecondaryVIPListFromContainer(cont)

	return list, err
}
