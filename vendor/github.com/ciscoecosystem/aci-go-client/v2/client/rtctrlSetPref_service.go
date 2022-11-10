package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateRtctrlSetPref(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetPrefAttr models.RtctrlSetPrefAttributes) (*models.RtctrlSetPref, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetPref)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetPref, tenant, action_rule_profile)
	rtctrlSetPref := models.NewRtctrlSetPref(rn, parentDn, description, nameAlias, rtctrlSetPrefAttr)
	err := sm.Save(rtctrlSetPref)
	return rtctrlSetPref, err
}

func (sm *ServiceManager) ReadRtctrlSetPref(action_rule_profile string, tenant string) (*models.RtctrlSetPref, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetPref, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetPref := models.RtctrlSetPrefFromContainer(cont)
	return rtctrlSetPref, nil
}

func (sm *ServiceManager) DeleteRtctrlSetPref(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetPref, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsetprefClassName)
}

func (sm *ServiceManager) UpdateRtctrlSetPref(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetPrefAttr models.RtctrlSetPrefAttributes) (*models.RtctrlSetPref, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetPref)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetPref, tenant, action_rule_profile)
	rtctrlSetPref := models.NewRtctrlSetPref(rn, parentDn, description, nameAlias, rtctrlSetPrefAttr)
	rtctrlSetPref.Status = "modified"
	err := sm.Save(rtctrlSetPref)
	return rtctrlSetPref, err
}

func (sm *ServiceManager) ListRtctrlSetPref(action_rule_profile string, tenant string) ([]*models.RtctrlSetPref, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetPref.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RtctrlSetPrefListFromContainer(cont)
	return list, err
}
