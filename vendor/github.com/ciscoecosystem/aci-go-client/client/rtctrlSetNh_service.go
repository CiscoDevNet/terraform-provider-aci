package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRtctrlSetNh(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetNhAttr models.RtctrlSetNhAttributes) (*models.RtctrlSetNh, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetNh)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetNh, tenant, action_rule_profile)
	rtctrlSetNh := models.NewRtctrlSetNh(rn, parentDn, description, nameAlias, rtctrlSetNhAttr)
	err := sm.Save(rtctrlSetNh)
	return rtctrlSetNh, err
}

func (sm *ServiceManager) ReadRtctrlSetNh(action_rule_profile string, tenant string) (*models.RtctrlSetNh, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetNh, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetNh := models.RtctrlSetNhFromContainer(cont)
	return rtctrlSetNh, nil
}

func (sm *ServiceManager) DeleteRtctrlSetNh(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetNh, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsetnhClassName)
}

func (sm *ServiceManager) UpdateRtctrlSetNh(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetNhAttr models.RtctrlSetNhAttributes) (*models.RtctrlSetNh, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetNh)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetNh, tenant, action_rule_profile)
	rtctrlSetNh := models.NewRtctrlSetNh(rn, parentDn, description, nameAlias, rtctrlSetNhAttr)
	rtctrlSetNh.Status = "modified"
	err := sm.Save(rtctrlSetNh)
	return rtctrlSetNh, err
}

func (sm *ServiceManager) ListRtctrlSetNh(action_rule_profile string, tenant string) ([]*models.RtctrlSetNh, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetNh.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RtctrlSetNhListFromContainer(cont)
	return list, err
}
