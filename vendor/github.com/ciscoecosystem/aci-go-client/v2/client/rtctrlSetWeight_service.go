package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateRtctrlSetWeight(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetWeightAttr models.RtctrlSetWeightAttributes) (*models.RtctrlSetWeight, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetWeight)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetWeight, tenant, action_rule_profile)
	rtctrlSetWeight := models.NewRtctrlSetWeight(rn, parentDn, description, nameAlias, rtctrlSetWeightAttr)
	err := sm.Save(rtctrlSetWeight)
	return rtctrlSetWeight, err
}

func (sm *ServiceManager) ReadRtctrlSetWeight(action_rule_profile string, tenant string) (*models.RtctrlSetWeight, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetWeight, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetWeight := models.RtctrlSetWeightFromContainer(cont)
	return rtctrlSetWeight, nil
}

func (sm *ServiceManager) DeleteRtctrlSetWeight(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetWeight, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsetweightClassName)
}

func (sm *ServiceManager) UpdateRtctrlSetWeight(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetWeightAttr models.RtctrlSetWeightAttributes) (*models.RtctrlSetWeight, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetWeight)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetWeight, tenant, action_rule_profile)
	rtctrlSetWeight := models.NewRtctrlSetWeight(rn, parentDn, description, nameAlias, rtctrlSetWeightAttr)
	rtctrlSetWeight.Status = "modified"
	err := sm.Save(rtctrlSetWeight)
	return rtctrlSetWeight, err
}

func (sm *ServiceManager) ListRtctrlSetWeight(action_rule_profile string, tenant string) ([]*models.RtctrlSetWeight, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetWeight.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RtctrlSetWeightListFromContainer(cont)
	return list, err
}
