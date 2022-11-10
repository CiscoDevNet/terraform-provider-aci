package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateRtctrlSetComm(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetCommAttr models.RtctrlSetCommAttributes) (*models.RtctrlSetComm, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetComm)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetComm, tenant, action_rule_profile)
	rtctrlSetComm := models.NewRtctrlSetComm(rn, parentDn, description, nameAlias, rtctrlSetCommAttr)
	err := sm.Save(rtctrlSetComm)
	return rtctrlSetComm, err
}

func (sm *ServiceManager) ReadRtctrlSetComm(action_rule_profile string, tenant string) (*models.RtctrlSetComm, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetComm, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetComm := models.RtctrlSetCommFromContainer(cont)
	return rtctrlSetComm, nil
}

func (sm *ServiceManager) DeleteRtctrlSetComm(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetComm, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsetcommClassName)
}

func (sm *ServiceManager) UpdateRtctrlSetComm(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetCommAttr models.RtctrlSetCommAttributes) (*models.RtctrlSetComm, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetComm)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetComm, tenant, action_rule_profile)
	rtctrlSetComm := models.NewRtctrlSetComm(rn, parentDn, description, nameAlias, rtctrlSetCommAttr)
	rtctrlSetComm.Status = "modified"
	err := sm.Save(rtctrlSetComm)
	return rtctrlSetComm, err
}

func (sm *ServiceManager) ListRtctrlSetComm(action_rule_profile string, tenant string) ([]*models.RtctrlSetComm, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetComm.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RtctrlSetCommListFromContainer(cont)
	return list, err
}
