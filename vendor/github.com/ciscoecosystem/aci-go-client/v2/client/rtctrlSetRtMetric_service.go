package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateRtctrlSetRtMetric(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetRtMetricAttr models.RtctrlSetRtMetricAttributes) (*models.RtctrlSetRtMetric, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetRtMetric)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetRtMetric, tenant, action_rule_profile)
	rtctrlSetRtMetric := models.NewRtctrlSetRtMetric(rn, parentDn, description, nameAlias, rtctrlSetRtMetricAttr)
	err := sm.Save(rtctrlSetRtMetric)
	return rtctrlSetRtMetric, err
}

func (sm *ServiceManager) ReadRtctrlSetRtMetric(action_rule_profile string, tenant string) (*models.RtctrlSetRtMetric, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetRtMetric, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetRtMetric := models.RtctrlSetRtMetricFromContainer(cont)
	return rtctrlSetRtMetric, nil
}

func (sm *ServiceManager) DeleteRtctrlSetRtMetric(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetRtMetric, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsetrtmetricClassName)
}

func (sm *ServiceManager) UpdateRtctrlSetRtMetric(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetRtMetricAttr models.RtctrlSetRtMetricAttributes) (*models.RtctrlSetRtMetric, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetRtMetric)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetRtMetric, tenant, action_rule_profile)
	rtctrlSetRtMetric := models.NewRtctrlSetRtMetric(rn, parentDn, description, nameAlias, rtctrlSetRtMetricAttr)
	rtctrlSetRtMetric.Status = "modified"
	err := sm.Save(rtctrlSetRtMetric)
	return rtctrlSetRtMetric, err
}

func (sm *ServiceManager) ListRtctrlSetRtMetric(action_rule_profile string, tenant string) ([]*models.RtctrlSetRtMetric, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetRtMetric.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RtctrlSetRtMetricListFromContainer(cont)
	return list, err
}
