package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateRtctrlSetDamp(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetDampAttr models.RtctrlSetDampAttributes) (*models.RtctrlSetDamp, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetDamp)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetDamp, tenant, action_rule_profile)
	rtctrlSetDamp := models.NewRtctrlSetDamp(rn, parentDn, description, nameAlias, rtctrlSetDampAttr)
	err := sm.Save(rtctrlSetDamp)
	return rtctrlSetDamp, err
}

func (sm *ServiceManager) ReadRtctrlSetDamp(action_rule_profile string, tenant string) (*models.RtctrlSetDamp, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetDamp, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetDamp := models.RtctrlSetDampFromContainer(cont)
	return rtctrlSetDamp, nil
}

func (sm *ServiceManager) DeleteRtctrlSetDamp(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetDamp, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsetdampClassName)
}

func (sm *ServiceManager) UpdateRtctrlSetDamp(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetDampAttr models.RtctrlSetDampAttributes) (*models.RtctrlSetDamp, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetDamp)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetDamp, tenant, action_rule_profile)
	rtctrlSetDamp := models.NewRtctrlSetDamp(rn, parentDn, description, nameAlias, rtctrlSetDampAttr)
	rtctrlSetDamp.Status = "modified"
	err := sm.Save(rtctrlSetDamp)
	return rtctrlSetDamp, err
}

func (sm *ServiceManager) ListRtctrlSetDamp(action_rule_profile string, tenant string) ([]*models.RtctrlSetDamp, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetDamp.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RtctrlSetDampListFromContainer(cont)
	return list, err
}
