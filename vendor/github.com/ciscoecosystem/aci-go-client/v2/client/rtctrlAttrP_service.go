package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateActionRuleProfile(name string, tenant string, description string, nameAlias string, rtctrlAttrPAttr models.ActionRuleProfileAttributes) (*models.ActionRuleProfile, error) {
	rn := fmt.Sprintf(models.RnrtctrlAttrP, name)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlAttrP, tenant)
	rtctrlAttrP := models.NewActionRuleProfile(rn, parentDn, description, nameAlias, rtctrlAttrPAttr)
	err := sm.Save(rtctrlAttrP)
	return rtctrlAttrP, err
}

func (sm *ServiceManager) ReadActionRuleProfile(name string, tenant string) (*models.ActionRuleProfile, error) {
	dn := fmt.Sprintf(models.DnrtctrlAttrP, tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlAttrP := models.ActionRuleProfileFromContainer(cont)
	return rtctrlAttrP, nil
}

func (sm *ServiceManager) DeleteActionRuleProfile(name string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlAttrP, tenant, name)
	return sm.DeleteByDn(dn, models.RtctrlattrpClassName)
}

func (sm *ServiceManager) UpdateActionRuleProfile(name string, tenant string, description string, nameAlias string, rtctrlAttrPAttr models.ActionRuleProfileAttributes) (*models.ActionRuleProfile, error) {
	rn := fmt.Sprintf(models.RnrtctrlAttrP, name)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlAttrP, tenant)
	rtctrlAttrP := models.NewActionRuleProfile(rn, parentDn, description, nameAlias, rtctrlAttrPAttr)

	rtctrlAttrP.Status = "modified"
	err := sm.Save(rtctrlAttrP)
	return rtctrlAttrP, err

}

func (sm *ServiceManager) ListActionRuleProfile(tenant string) ([]*models.ActionRuleProfile, error) {

	dnUrl := fmt.Sprintf("%s/uni/tn-%s/rtctrlAttrP.json", models.BaseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ActionRuleProfileListFromContainer(cont)

	return list, err
}
