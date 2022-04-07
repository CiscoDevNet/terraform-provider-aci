package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRtctrlSetAddComm(community string, action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetAddCommAttr models.RtctrlSetAddCommAttributes) (*models.RtctrlSetAddComm, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetAddComm, community)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetAddComm, tenant, action_rule_profile)
	rtctrlSetAddComm := models.NewRtctrlSetAddComm(rn, parentDn, description, nameAlias, rtctrlSetAddCommAttr)
	err := sm.Save(rtctrlSetAddComm)
	return rtctrlSetAddComm, err
}

func (sm *ServiceManager) ReadRtctrlSetAddComm(community string, action_rule_profile string, tenant string) (*models.RtctrlSetAddComm, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetAddComm, tenant, action_rule_profile, community)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetAddComm := models.RtctrlSetAddCommFromContainer(cont)
	return rtctrlSetAddComm, nil
}

func (sm *ServiceManager) DeleteRtctrlSetAddComm(community string, action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetAddComm, tenant, action_rule_profile, community)
	return sm.DeleteByDn(dn, models.RtctrlsetaddcommClassName)
}

func (sm *ServiceManager) UpdateRtctrlSetAddComm(community string, action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetAddCommAttr models.RtctrlSetAddCommAttributes) (*models.RtctrlSetAddComm, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetAddComm, community)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetAddComm, tenant, action_rule_profile)
	rtctrlSetAddComm := models.NewRtctrlSetAddComm(rn, parentDn, description, nameAlias, rtctrlSetAddCommAttr)
	rtctrlSetAddComm.Status = "modified"
	err := sm.Save(rtctrlSetAddComm)
	return rtctrlSetAddComm, err
}

func (sm *ServiceManager) ListRtctrlSetAddComm(action_rule_profile string, tenant string) ([]*models.RtctrlSetAddComm, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetAddComm.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RtctrlSetAddCommListFromContainer(cont)
	return list, err
}
