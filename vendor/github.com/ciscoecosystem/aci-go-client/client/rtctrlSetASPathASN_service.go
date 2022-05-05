package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateASNumber(order string, set_as_path_criteria string, action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetASPathASNAttr models.ASNumberAttributes) (*models.ASNumber, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetASPathASN, order)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetASPathASN, tenant, action_rule_profile, set_as_path_criteria)
	rtctrlSetASPathASN := models.NewASNumber(rn, parentDn, description, nameAlias, rtctrlSetASPathASNAttr)
	err := sm.Save(rtctrlSetASPathASN)
	return rtctrlSetASPathASN, err
}

func (sm *ServiceManager) ReadASNumber(order string, set_as_path_criteria string, action_rule_profile string, tenant string) (*models.ASNumber, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetASPathASN, tenant, action_rule_profile, set_as_path_criteria, order)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetASPathASN := models.ASNumberFromContainer(cont)
	return rtctrlSetASPathASN, nil
}

func (sm *ServiceManager) DeleteASNumber(order string, set_as_path_criteria string, action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetASPathASN, tenant, action_rule_profile, set_as_path_criteria, order)
	return sm.DeleteByDn(dn, models.RtctrlsetaspathasnClassName)
}

func (sm *ServiceManager) UpdateASNumber(order string, set_as_path_criteria string, action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetASPathASNAttr models.ASNumberAttributes) (*models.ASNumber, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetASPathASN, order)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetASPathASN, tenant, action_rule_profile, set_as_path_criteria)
	rtctrlSetASPathASN := models.NewASNumber(rn, parentDn, description, nameAlias, rtctrlSetASPathASNAttr)
	rtctrlSetASPathASN.Status = "modified"
	err := sm.Save(rtctrlSetASPathASN)
	return rtctrlSetASPathASN, err
}

func (sm *ServiceManager) ListASNumber(set_as_path_criteria string, action_rule_profile string, tenant string) ([]*models.ASNumber, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/saspath-%s/rtctrlSetASPathASN.json", models.BaseurlStr, tenant, action_rule_profile, set_as_path_criteria)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ASNumberListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) ListSetAsPathASNs(parentDn string) ([]*models.ASNumber, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "rtctrlSetASPathASN")
	cont, err := sm.GetViaURL(dnUrl)

	if err != nil {
		return nil, err
	} else {
		contList := models.ASNumberListFromContainer(cont)
		return contList, nil
	}
}
