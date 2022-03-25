package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateMatchRuleBasedonCommunityRegularExpression(commType string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchCommRegexTermAttr models.MatchRuleBasedonCommunityRegularExpressionAttributes) (*models.MatchRuleBasedonCommunityRegularExpression, error) {
	rn := fmt.Sprintf(models.RnrtctrlMatchCommRegexTerm, commType)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlMatchCommRegexTerm, tenant, match_rule)
	rtctrlMatchCommRegexTerm := models.NewMatchRuleBasedonCommunityRegularExpression(rn, parentDn, description, nameAlias, rtctrlMatchCommRegexTermAttr)
	err := sm.Save(rtctrlMatchCommRegexTerm)
	return rtctrlMatchCommRegexTerm, err
}

func (sm *ServiceManager) ReadMatchRuleBasedonCommunityRegularExpression(commType string, match_rule string, tenant string) (*models.MatchRuleBasedonCommunityRegularExpression, error) {
	dn := fmt.Sprintf(models.DnrtctrlMatchCommRegexTerm, tenant, match_rule, commType)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlMatchCommRegexTerm := models.MatchRuleBasedonCommunityRegularExpressionFromContainer(cont)
	return rtctrlMatchCommRegexTerm, nil
}

func (sm *ServiceManager) DeleteMatchRuleBasedonCommunityRegularExpression(commType string, match_rule string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlMatchCommRegexTerm, tenant, match_rule, commType)
	return sm.DeleteByDn(dn, models.RtctrlmatchcommregextermClassName)
}

func (sm *ServiceManager) UpdateMatchRuleBasedonCommunityRegularExpression(commType string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchCommRegexTermAttr models.MatchRuleBasedonCommunityRegularExpressionAttributes) (*models.MatchRuleBasedonCommunityRegularExpression, error) {
	rn := fmt.Sprintf(models.RnrtctrlMatchCommRegexTerm, commType)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlMatchCommRegexTerm, tenant, match_rule)
	rtctrlMatchCommRegexTerm := models.NewMatchRuleBasedonCommunityRegularExpression(rn, parentDn, description, nameAlias, rtctrlMatchCommRegexTermAttr)
	rtctrlMatchCommRegexTerm.Status = "modified"
	err := sm.Save(rtctrlMatchCommRegexTerm)
	return rtctrlMatchCommRegexTerm, err
}

func (sm *ServiceManager) ListMatchRuleBasedonCommunityRegularExpression(match_rule string, tenant string) ([]*models.MatchRuleBasedonCommunityRegularExpression, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/subj-%s/rtctrlMatchCommRegexTerm.json", models.BaseurlStr, tenant, match_rule)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MatchRuleBasedonCommunityRegularExpressionListFromContainer(cont)
	return list, err
}
