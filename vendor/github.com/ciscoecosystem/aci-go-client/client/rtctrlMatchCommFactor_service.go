package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateMatchCommunityFactor(community string, match_community_term string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchCommFactorAttr models.MatchCommunityFactorAttributes) (*models.MatchCommunityFactor, error) {
	rn := fmt.Sprintf(models.RnrtctrlMatchCommFactor, community)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlMatchCommFactor, tenant, match_rule, match_community_term)
	rtctrlMatchCommFactor := models.NewMatchCommunityFactor(rn, parentDn, description, nameAlias, rtctrlMatchCommFactorAttr)
	err := sm.Save(rtctrlMatchCommFactor)
	return rtctrlMatchCommFactor, err
}

func (sm *ServiceManager) ReadMatchCommunityFactor(community string, match_community_term string, match_rule string, tenant string) (*models.MatchCommunityFactor, error) {
	dn := fmt.Sprintf(models.DnrtctrlMatchCommFactor, tenant, match_rule, match_community_term, community)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlMatchCommFactor := models.MatchCommunityFactorFromContainer(cont)
	return rtctrlMatchCommFactor, nil
}

func (sm *ServiceManager) DeleteMatchCommunityFactor(community string, match_community_term string, match_rule string, tenant string) error {
	dn := fmt.Sprintf(models.DnrtctrlMatchCommFactor, tenant, match_rule, match_community_term, community)
	return sm.DeleteByDn(dn, models.RtctrlmatchcommfactorClassName)
}

func (sm *ServiceManager) UpdateMatchCommunityFactor(community string, match_community_term string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchCommFactorAttr models.MatchCommunityFactorAttributes) (*models.MatchCommunityFactor, error) {
	rn := fmt.Sprintf(models.RnrtctrlMatchCommFactor, community)
	parentDn := fmt.Sprintf(models.ParentDnrtctrlMatchCommFactor, tenant, match_rule, match_community_term)
	rtctrlMatchCommFactor := models.NewMatchCommunityFactor(rn, parentDn, description, nameAlias, rtctrlMatchCommFactorAttr)
	rtctrlMatchCommFactor.Status = "modified"
	err := sm.Save(rtctrlMatchCommFactor)
	return rtctrlMatchCommFactor, err
}

func (sm *ServiceManager) ListMatchCommunityFactor(match_community_term string, match_rule string, tenant string) ([]*models.MatchCommunityFactor, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/subj-%s/commtrm-%s/rtctrlMatchCommFactor.json", models.BaseurlStr, tenant, match_rule, match_community_term)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MatchCommunityFactorListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) ListMatchCommFactorsFromCommunityTerm(parentDn string) ([]*models.MatchCommunityFactor, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "rtctrlMatchCommFactor")
	cont, err := sm.GetViaURL(dnUrl)

	rtctrlMatchCommFactorsData := models.MatchCommunityFactorListFromContainer(cont)

	return rtctrlMatchCommFactorsData, err

}
