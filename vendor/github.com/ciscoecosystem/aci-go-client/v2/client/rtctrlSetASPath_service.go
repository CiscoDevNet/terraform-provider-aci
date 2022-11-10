package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateSetASPath(criteria string, action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetASPathAttr models.SetASPathAttributes) (*models.SetASPath, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetASPath, criteria)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetASPath, tenant, action_rule_profile)
	rtctrlSetASPath := models.NewSetASPath(rn, parentDn, description, nameAlias, rtctrlSetASPathAttr)
	err := sm.Save(rtctrlSetASPath)
	return rtctrlSetASPath, err
}

func (sm *ServiceManager) ReadSetASPath(criteria string, action_rule_profile string, tenant string) (*models.SetASPath, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetASPath, tenant, action_rule_profile, criteria)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetASPath := models.SetASPathFromContainer(cont)
	return rtctrlSetASPath, nil
}

func (sm *ServiceManager) DeleteSetASPath(criteria string, action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetASPath, tenant, action_rule_profile, criteria)
	return sm.DeleteByDn(dn, models.RtctrlsetaspathClassName)
}

func (sm *ServiceManager) UpdateSetASPath(criteria string, action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetASPathAttr models.SetASPathAttributes) (*models.SetASPath, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetASPath, criteria)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetASPath, tenant, action_rule_profile)
	rtctrlSetASPath := models.NewSetASPath(rn, parentDn, description, nameAlias, rtctrlSetASPathAttr)
	rtctrlSetASPath.Status = "modified"
	err := sm.Save(rtctrlSetASPath)
	return rtctrlSetASPath, err
}

func (sm *ServiceManager) ListSetASPath(action_rule_profile string, tenant string) ([]*models.SetASPath, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetASPath.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SetASPathListFromContainer(cont)
	return list, err
}
