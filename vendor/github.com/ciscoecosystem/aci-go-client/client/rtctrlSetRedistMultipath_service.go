package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRedistributeMultipathAction(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetRedistMultipathAttr models.RedistributeMultipathActionAttributes) (*models.RedistributeMultipathAction, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetRedistMultipath)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetRedistMultipath, tenant, action_rule_profile)
	rtctrlSetRedistMultipath := models.NewRedistributeMultipathAction(rn, parentDn, description, nameAlias, rtctrlSetRedistMultipathAttr)
	err := sm.Save(rtctrlSetRedistMultipath)
	return rtctrlSetRedistMultipath, err
}

func (sm *ServiceManager) ReadRedistributeMultipathAction(action_rule_profile string, tenant string) (*models.RedistributeMultipathAction, error) {
	dn := fmt.Sprintf(models.DnrtctrlSetRedistMultipath, tenant, action_rule_profile)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlSetRedistMultipath := models.RedistributeMultipathActionFromContainer(cont)
	return rtctrlSetRedistMultipath, nil
}

func (sm *ServiceManager) DeleteRedistributeMultipathAction(action_rule_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSetRedistMultipath, tenant, action_rule_profile)
	return sm.DeleteByDn(dn, models.RtctrlsetredistmultipathClassName)
}

func (sm *ServiceManager) UpdateRedistributeMultipathAction(action_rule_profile string, tenant string, description string, nameAlias string, rtctrlSetRedistMultipathAttr models.RedistributeMultipathActionAttributes) (*models.RedistributeMultipathAction, error) {
	rn := fmt.Sprintf(models.RnrtctrlSetRedistMultipath)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSetRedistMultipath, tenant, action_rule_profile)
	rtctrlSetRedistMultipath := models.NewRedistributeMultipathAction(rn, parentDn, description, nameAlias, rtctrlSetRedistMultipathAttr)
	rtctrlSetRedistMultipath.Status = "modified"
	err := sm.Save(rtctrlSetRedistMultipath)
	return rtctrlSetRedistMultipath, err
}

func (sm *ServiceManager) ListRedistributeMultipathAction(action_rule_profile string, tenant string) ([]*models.RedistributeMultipathAction, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/attr-%s/rtctrlSetRedistMultipath.json", models.BaseurlStr, tenant, action_rule_profile)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RedistributeMultipathActionListFromContainer(cont)
	return list, err
}
