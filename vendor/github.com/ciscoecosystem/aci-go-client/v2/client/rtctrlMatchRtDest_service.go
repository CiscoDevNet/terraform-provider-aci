package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateMatchRouteDestinationRule(ip string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchRtDestAttr models.MatchRouteDestinationRuleAttributes) (*models.MatchRouteDestinationRule, error) {
	rn := fmt.Sprintf(models.RnrtctrlMatchRtDest, ip)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlMatchRtDest, tenant, match_rule)
	rtctrlMatchRtDest := models.NewMatchRouteDestinationRule(rn, parentDn, description, nameAlias, rtctrlMatchRtDestAttr)
	err := sm.Save(rtctrlMatchRtDest)
	return rtctrlMatchRtDest, err
}

func (sm *ServiceManager) ReadMatchRouteDestinationRule(ip string, match_rule string, tenant string) (*models.MatchRouteDestinationRule, error) {
	dn := fmt.Sprintf(models.DnrtctrlMatchRtDest, tenant, match_rule, ip)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlMatchRtDest := models.MatchRouteDestinationRuleFromContainer(cont)
	return rtctrlMatchRtDest, nil
}

func (sm *ServiceManager) DeleteMatchRouteDestinationRule(ip string, match_rule string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlMatchRtDest, tenant, match_rule, ip)
	return sm.DeleteByDn(dn, models.RtctrlmatchrtdestClassName)
}

func (sm *ServiceManager) UpdateMatchRouteDestinationRule(ip string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchRtDestAttr models.MatchRouteDestinationRuleAttributes) (*models.MatchRouteDestinationRule, error) {
	rn := fmt.Sprintf(models.RnrtctrlMatchRtDest, ip)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlMatchRtDest, tenant, match_rule)
	rtctrlMatchRtDest := models.NewMatchRouteDestinationRule(rn, parentDn, description, nameAlias, rtctrlMatchRtDestAttr)
	rtctrlMatchRtDest.Status = "modified"
	err := sm.Save(rtctrlMatchRtDest)
	return rtctrlMatchRtDest, err
}

func (sm *ServiceManager) ListMatchRouteDestinationRule(match_rule string, tenant string) ([]*models.MatchRouteDestinationRule, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/subj-%s/rtctrlMatchRtDest.json", models.BaseurlStr, tenant, match_rule)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MatchRouteDestinationRuleListFromContainer(cont)
	return list, err
}
