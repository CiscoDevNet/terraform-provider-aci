package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateMatchCommunityTerm(name string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchCommTermAttr models.MatchCommunityTermAttributes) (*models.MatchCommunityTerm, error) {
	rn := fmt.Sprintf(models.RnrtctrlMatchCommTerm, name)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlMatchCommTerm, tenant, match_rule)
	rtctrlMatchCommTerm := models.NewMatchCommunityTerm(rn, parentDn, description, nameAlias, rtctrlMatchCommTermAttr)
	err := sm.Save(rtctrlMatchCommTerm)
	return rtctrlMatchCommTerm, err
}

func (sm *ServiceManager) ReadMatchCommunityTerm(name string, match_rule string, tenant string) (*models.MatchCommunityTerm, error) {
	dn := fmt.Sprintf(models.DnrtctrlMatchCommTerm, tenant, match_rule, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlMatchCommTerm := models.MatchCommunityTermFromContainer(cont)
	return rtctrlMatchCommTerm, nil
}

func (sm *ServiceManager) DeleteMatchCommunityTerm(name string, match_rule string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlMatchCommTerm, tenant, match_rule, name)
	return sm.DeleteByDn(dn, models.RtctrlmatchcommtermClassName)
}

func (sm *ServiceManager) UpdateMatchCommunityTerm(name string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchCommTermAttr models.MatchCommunityTermAttributes) (*models.MatchCommunityTerm, error) {
	rn := fmt.Sprintf(models.RnrtctrlMatchCommTerm, name)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlMatchCommTerm, tenant, match_rule)
	rtctrlMatchCommTerm := models.NewMatchCommunityTerm(rn, parentDn, description, nameAlias, rtctrlMatchCommTermAttr)
	rtctrlMatchCommTerm.Status = "modified"
	err := sm.Save(rtctrlMatchCommTerm)
	return rtctrlMatchCommTerm, err
}

func (sm *ServiceManager) ListMatchCommunityTerm(match_rule string, tenant string) ([]*models.MatchCommunityTerm, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/subj-%s/rtctrlMatchCommTerm.json", models.BaseurlStr, tenant, match_rule)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MatchCommunityTermListFromContainer(cont)
	return list, err
}
