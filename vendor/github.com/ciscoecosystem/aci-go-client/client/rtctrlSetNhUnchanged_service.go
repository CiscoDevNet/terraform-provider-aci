package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateNexthopUnchangedAction(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetNhUnchangedAttr models.NexthopUnchangedActionAttributes) (*models.NexthopUnchangedAction, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetNhUnchanged)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetNhUnchanged, tenant, action_rule_profile)
	rtctrlSetNhUnchanged := models.NewNexthopUnchangedAction(rn, parentDn, description, nameAlias, rtctrlSetNhUnchangedAttr)
	err := sm.Save(rtctrlSetNhUnchanged)
	return rtctrlSetNhUnchanged, err
}

func (sm *ServiceManager) ReadNexthopUnchangedAction(action_rule_profile string, tenant string) (*models.NexthopUnchangedAction, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetNhUnchanged, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetNhUnchanged := models.NexthopUnchangedActionFromContainer(cont)
	return rtctrlSetNhUnchanged, nil
}

func (sm *ServiceManager) DeleteNexthopUnchangedAction(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetNhUnchanged, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsetnhunchangedClassName)
}

func (sm *ServiceManager) UpdateNexthopUnchangedAction(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetNhUnchangedAttr models.NexthopUnchangedActionAttributes) (*models.NexthopUnchangedAction, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetNhUnchanged)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetNhUnchanged, tenant, action_rule_profile)
	rtctrlSetNhUnchanged := models.NewNexthopUnchangedAction(rn, parentDn, description, nameAlias, rtctrlSetNhUnchangedAttr)
	rtctrlSetNhUnchanged.Status = "modified"
	err := sm.Save(rtctrlSetNhUnchanged)
	return rtctrlSetNhUnchanged, err
}

func (sm *ServiceManager) ListNexthopUnchangedAction(action_rule_profile string, tenant string) ([]*models.NexthopUnchangedAction, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetNhUnchanged.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.NexthopUnchangedActionListFromContainer(cont)
	return list, err
}
