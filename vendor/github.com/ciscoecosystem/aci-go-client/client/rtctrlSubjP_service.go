package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateMatchRule(name string, tenant string, description string, nameAlias string, rtctrlSubjPAttr models.MatchRuleAttributes) (*models.MatchRule, error) {
	rn := fmt.Sprintf(models.RnrtctrlSubjP, name)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSubjP, tenant)
	rtctrlSubjP := models.NewMatchRule(rn, parentDn, description, nameAlias, rtctrlSubjPAttr)
	err := sm.Save(rtctrlSubjP)
	return rtctrlSubjP, err
}

func (sm *ServiceManager) ReadMatchRule(name string, tenant string) (*models.MatchRule, error) {
	dn := fmt.Sprintf(models.DnrtctrlSubjP, tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSubjP := models.MatchRuleFromContainer(cont)
	return rtctrlSubjP, nil
}

func (sm *ServiceManager) DeleteMatchRule(name string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlSubjP, tenant, name)
	return sm.DeleteByDn(dn, models.RtctrlsubjpClassName)
}

func (sm *ServiceManager) UpdateMatchRule(name string, tenant string, description string, nameAlias string, rtctrlSubjPAttr models.MatchRuleAttributes) (*models.MatchRule, error) {
	rn := fmt.Sprintf(models.RnrtctrlSubjP, name)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlSubjP, tenant)
	rtctrlSubjP := models.NewMatchRule(rn, parentDn, description, nameAlias, rtctrlSubjPAttr)
	rtctrlSubjP.Status = "modified"
	err := sm.Save(rtctrlSubjP)
	return rtctrlSubjP, err
}

func (sm *ServiceManager) ListMatchRule(tenant string) ([]*models.MatchRule, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/rtctrlSubjP.json", models.BaseurlStr, tenant)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MatchRuleListFromContainer(cont)
	return list, err
}
