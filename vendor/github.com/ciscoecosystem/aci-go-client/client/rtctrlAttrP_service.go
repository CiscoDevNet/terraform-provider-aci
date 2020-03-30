package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateActionRuleProfile(name string, tenant string, description string, rtctrlAttrPattr models.ActionRuleProfileAttributes) (*models.ActionRuleProfile, error) {
	rn := fmt.Sprintf("attr-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	rtctrlAttrP := models.NewActionRuleProfile(rn, parentDn, description, rtctrlAttrPattr)
	err := sm.Save(rtctrlAttrP)
	return rtctrlAttrP, err
}

func (sm *ServiceManager) ReadActionRuleProfile(name string, tenant string) (*models.ActionRuleProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/attr-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlAttrP := models.ActionRuleProfileFromContainer(cont)
	return rtctrlAttrP, nil
}

func (sm *ServiceManager) DeleteActionRuleProfile(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/attr-%s", tenant, name)
	return sm.DeleteByDn(dn, models.RtctrlattrpClassName)
}

func (sm *ServiceManager) UpdateActionRuleProfile(name string, tenant string, description string, rtctrlAttrPattr models.ActionRuleProfileAttributes) (*models.ActionRuleProfile, error) {
	rn := fmt.Sprintf("attr-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	rtctrlAttrP := models.NewActionRuleProfile(rn, parentDn, description, rtctrlAttrPattr)

	rtctrlAttrP.Status = "modified"
	err := sm.Save(rtctrlAttrP)
	return rtctrlAttrP, err

}

func (sm *ServiceManager) ListActionRuleProfile(tenant string) ([]*models.ActionRuleProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/rtctrlAttrP.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ActionRuleProfileListFromContainer(cont)

	return list, err
}
