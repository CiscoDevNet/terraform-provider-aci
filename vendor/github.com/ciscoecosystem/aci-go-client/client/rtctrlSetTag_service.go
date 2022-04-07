package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRtctrlSetTag(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetTagAttr models.RtctrlSetTagAttributes) (*models.RtctrlSetTag, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetTag)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetTag, tenant, action_rule_profile)
	rtctrlSetTag := models.NewRtctrlSetTag(rn, parentDn, description, nameAlias, rtctrlSetTagAttr)
	err := sm.Save(rtctrlSetTag)
	return rtctrlSetTag, err
}

func (sm *ServiceManager) ReadRtctrlSetTag(action_rule_profile string, tenant string) (*models.RtctrlSetTag, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetTag, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetTag := models.RtctrlSetTagFromContainer(cont)
	return rtctrlSetTag, nil
}

func (sm *ServiceManager) DeleteRtctrlSetTag(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetTag, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsettagClassName)
}

func (sm *ServiceManager) UpdateRtctrlSetTag(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetTagAttr models.RtctrlSetTagAttributes) (*models.RtctrlSetTag, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetTag)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetTag, tenant, action_rule_profile)
	rtctrlSetTag := models.NewRtctrlSetTag(rn, parentDn, description, nameAlias, rtctrlSetTagAttr)
	rtctrlSetTag.Status = "modified"
	err := sm.Save(rtctrlSetTag)
	return rtctrlSetTag, err
}

func (sm *ServiceManager) ListRtctrlSetTag(action_rule_profile string, tenant string) ([]*models.RtctrlSetTag, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetTag.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RtctrlSetTagListFromContainer(cont)
	return list, err
}
