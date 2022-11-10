package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateRtctrlSetRtMetricType(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetRtMetricTypeAttr models.RtctrlSetRtMetricTypeAttributes) (*models.RtctrlSetRtMetricType, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetRtMetricType)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetRtMetricType, tenant, action_rule_profile)
	rtctrlSetRtMetricType := models.NewRtctrlSetRtMetricType(rn, parentDn, description, nameAlias, rtctrlSetRtMetricTypeAttr)
	err := sm.Save(rtctrlSetRtMetricType)
	return rtctrlSetRtMetricType, err
}

func (sm *ServiceManager) ReadRtctrlSetRtMetricType(action_rule_profile string, tenant string) (*models.RtctrlSetRtMetricType, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetRtMetricType, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetRtMetricType := models.RtctrlSetRtMetricTypeFromContainer(cont)
	return rtctrlSetRtMetricType, nil
}

func (sm *ServiceManager) DeleteRtctrlSetRtMetricType(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetRtMetricType, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsetrtmetrictypeClassName)
}

func (sm *ServiceManager) UpdateRtctrlSetRtMetricType(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetRtMetricTypeAttr models.RtctrlSetRtMetricTypeAttributes) (*models.RtctrlSetRtMetricType, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetRtMetricType)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetRtMetricType, tenant, action_rule_profile)
	rtctrlSetRtMetricType := models.NewRtctrlSetRtMetricType(rn, parentDn, description, nameAlias, rtctrlSetRtMetricTypeAttr)
	rtctrlSetRtMetricType.Status = "modified"
	err := sm.Save(rtctrlSetRtMetricType)
	return rtctrlSetRtMetricType, err
}

func (sm *ServiceManager) ListRtctrlSetRtMetricType(action_rule_profile string, tenant string) ([]*models.RtctrlSetRtMetricType, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetRtMetricType.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RtctrlSetRtMetricTypeListFromContainer(cont)
	return list, err
}
